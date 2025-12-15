package provider

import (
	"context"
	"sync"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/pkg/logger"
)

type AirlineProvider struct {
	airAsia         AirlineInterface
	batikAir        AirlineInterface
	garudaIndonesia AirlineInterface
	lionAir         AirlineInterface
}

func NewAirlineProvider() *AirlineProvider {
	airAsia := NewAirAsiaProvider("./internal/provider/mock")
	batikAir := NewBatikAirProvider("./internal/provider/mock")
	garudaIndonesia := NewGarudaIndonesiaProvider("./internal/provider/mock")
	lionAir := NewLionAirProvider("./internal/provider/mock")
	return &AirlineProvider{
		airAsia:         airAsia,
		batikAir:        batikAir,
		garudaIndonesia: garudaIndonesia,
		lionAir:         lionAir,
	}
}

func (ap *AirlineProvider) SearchFlights(ctx context.Context, input domain.SearchRequest) (*domain.SearchResponse, error) {
	// Implementation for searching flights from this specific airline provider
	MAX_RETRY := 3
	TIMEOUT := 2 * time.Second

	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()

	type result struct {
		provider string
		flights  []domain.FlightInfo
		err      error
	}

	providers := []struct {
		name string
		fn   func() ([]domain.FlightInfo, error)
	}{
		{"airasia", func() ([]domain.FlightInfo, error) {
			return ap.airAsia.SearchFlights(ctx, input)
		}},
		{"batik", func() ([]domain.FlightInfo, error) {
			return ap.batikAir.SearchFlights(ctx, input)
		}},
		{"garuda", func() ([]domain.FlightInfo, error) {
			return ap.garudaIndonesia.SearchFlights(ctx, input)
		}},
		{"lion", func() ([]domain.FlightInfo, error) {
			return ap.lionAir.SearchFlights(ctx, input)
		}},
	}

	var wg sync.WaitGroup
	wg.Add(len(providers))

	ch := make(chan result, len(providers))
	for idx, p := range providers {
		// capture range variable
		provider := p

		go func(i int) {
			defer wg.Done()

			var (
				flights []domain.FlightInfo
				err     error
			)
			for range MAX_RETRY {
				if ctx.Err() != nil {
					err = ctx.Err()
					break
				}
				flights, err = provider.fn()
				// retry if error
				if err == nil {
					break
				}
			}
			ch <- result{provider: provider.name, flights: flights, err: err}
		}(idx)
	}

	wg.Wait()
	close(ch)

	resp := &domain.SearchResponse{
		Flights: []domain.FlightInfo{},
	}

	for res := range ch {
		// logger.InfoContext(ctx, "Provider Result: ", "flights", res.flights, "err", res.err)
		resp.Metadata.ProvidersQueried++
		if res.err != nil {
			logger.ErrorContext(ctx, "Provider failed", "err", res.err)
			resp.Metadata.ProvidersFailed++
			continue
		}
		logger.InfoContext(ctx, "Provider succeeded", "provider", res.provider, "len_flights", len(res.flights))
		resp.Metadata.TotalResults += len(res.flights)
		resp.Flights = append(resp.Flights, res.flights...)
		resp.Metadata.ProvidersSucceeded++
	}

	logger.InfoContext(ctx, "Total results", "total", len(resp.Flights))

	return resp, nil
}
