package service

import (
	"context"
	"sort"
	"time"

	"github.com/azcov/bookcabin_test/internal/config"
	"github.com/azcov/bookcabin_test/internal/consts"
	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/provider"
	"github.com/azcov/bookcabin_test/pkg/cache"
	"github.com/azcov/bookcabin_test/pkg/logger"
)

type FlightInterface interface {
	SerchFlight(ctx context.Context, input *domain.SearchRequest) (*domain.SearchResponse, error)
}

type flightService struct {
	airlaneProvider provider.AirlineAggregator
	cache           cache.Cache
}

func NewFlightService(cfg config.Config) FlightInterface {
	airlaneProvider := provider.NewAirlineProvider()
	return &flightService{
		airlaneProvider: airlaneProvider,
		cache:           cache.NewGoCache(cfg.Cache),
	}
}

func (fs *flightService) SerchFlight(ctx context.Context, input *domain.SearchRequest) (*domain.SearchResponse, error) {
	start := time.Now()

	// 1. Check Cache
	cacheKey := input.ToCacheKey()
	if cachedData, err := fs.cache.Get(cacheKey); err == nil {
		data := cachedData.(domain.SearchResponse)
		data.Metadata.CacheHit = true
		data.Metadata.SearchTimeMs = int(time.Since(start).Milliseconds())
		return &data, nil
	}

	// 2. Call Providers
	result, err := fs.airlaneProvider.SearchFlights(ctx, *input)
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}
	// logger.InfoContext(ctx, "Aggregated Flights: ", "count", len(result.Flights), "data", result.Flights)

	// 3. Filter Results
	result.Flights = fs.filterFlights(result.Flights, input.Filters)

	// 4. Calculate Best Value Score (Ranking)
	fs.calculateBestValue(result.Flights)

	// 5. Sort Results
	fs.sortFlights(result.Flights, input.Sort)

	// Update Metadata
	result.Metadata.TotalResults = len(result.Flights)
	result.SearchCriteria = *input
	result.Metadata.SearchTimeMs = int(time.Since(start).Milliseconds())
	result.Metadata.CacheHit = false

	// 6. Save to Cache
	fs.cache.Set(cacheKey, result)

	return result, nil
}

// --- Aggregation Logic ---

func (fs *flightService) filterFlights(flights []domain.FlightInfo, filters []domain.SearchFilter) []domain.FlightInfo {
	if len(filters) == 0 {
		return []domain.FlightInfo{}
	}

	filtered := make([]domain.FlightInfo, 0, len(flights))

	for _, f := range flights {
		keep := true
		for _, filter := range filters {
			if !fs.applyFilter(f, filter) {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func (fs *flightService) calculateBestValue(flights []domain.FlightInfo) {
	for _, f := range flights {
		f.CalculateBestValueScore()
	}
}

func (fs *flightService) applyFilter(f domain.FlightInfo, filter domain.SearchFilter) bool {
	switch filter.Key {
	case consts.FilterKeyMaxPrice:
		if val, ok := filter.Value.(float64); ok {
			return float64(f.Price.Amount) <= val
		}
	case consts.FilterKeyMinPrice:
		if val, ok := filter.Value.(float64); ok {
			return float64(f.Price.Amount) >= val
		}
	case consts.FilterKeyMaxStops:
		if val, ok := filter.Value.(float64); ok {
			return float64(f.Stops) <= val
		}
	case consts.FilterKeyMaxDuration: // Minutes
		if val, ok := filter.Value.(float64); ok {
			return float64(f.Duration.TotalMinutes) <= val
		}
	case consts.FilterKeyAirlines:
		if val, ok := filter.Value.(string); ok {
			return f.Airline.Code == val || f.Airline.Name == val
		}
	}
	return true
}

func (fs *flightService) sortFlights(flights []domain.FlightInfo, sortOpt domain.SortOption) {
	sort.SliceStable(flights, func(i, j int) bool {
		a, b := flights[i], flights[j]

		switch sortOpt.Key {
		case consts.SortKeyPrice:
			if sortOpt.Order == consts.SortOrderDesc {
				return a.Price.Amount > b.Price.Amount
			}
			return a.Price.Amount < b.Price.Amount

		case consts.SortKeyDuration:
			if sortOpt.Order == consts.SortOrderDesc {
				return a.Duration.TotalMinutes > b.Duration.TotalMinutes
			}
			return a.Duration.TotalMinutes < b.Duration.TotalMinutes

		case consts.SortKeyDepartureTime:
			if sortOpt.Order == consts.SortOrderDesc {
				return a.Departure.Timestamp > b.Departure.Timestamp
			}
			return a.Departure.Timestamp < b.Departure.Timestamp

		case consts.SortKeyArrivalTime:
			if sortOpt.Order == consts.SortOrderDesc {
				return a.Arrival.Timestamp > b.Arrival.Timestamp
			}
			return a.Arrival.Timestamp < b.Arrival.Timestamp

		default: // Default "Best Value" ranking
			if sortOpt.Order == consts.SortOrderDesc {
				return a.BestValueScore > b.BestValueScore
			}
			return a.BestValueScore < b.BestValueScore
		}
	})
}
