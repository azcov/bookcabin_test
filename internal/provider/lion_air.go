package provider

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/errors"
	lionair "github.com/azcov/bookcabin_test/internal/provider/lion_air"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/azcov/bookcabin_test/pkg/ratelimit"
)

type lionAirProvider struct {
	fileDir string
	rl      ratelimit.Limiter
}

func NewLionAirProvider(fileDir string) AirlineInterface {
	return &lionAirProvider{
		fileDir: fileDir,
		rl:      ratelimit.NewWithDuration(100, time.Second),
	}
}

func (ap *lionAirProvider) SearchFlights(ctx context.Context, input domain.SearchRequest) (flights []domain.FlightInfo, err error) {
	if !ap.rl.Allow() {
		return nil, errors.ErrLionAirRateLimitExceeded
	}
	// Implementation for searching flights from Batik Air
	start := time.Now()
	minDelay := 100
	maxDelay := 200

	defer func() {
		elapsed := time.Since(start)
		delay := util.RandomDuration(minDelay, maxDelay)
		wait := delay - elapsed
		if wait > 0 {
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				err = ctx.Err()
			}
		}
	}()

	flights, err = ap.callSearch(ctx, input)
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	return flights, nil
}

// callSearch simulates calling the Batik Air search API and returns filtered mock data.
func (ap *lionAirProvider) callSearch(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	flights := []domain.FlightInfo{}
	// Load mock file
	data, err := os.ReadFile(ap.fileDir + "/lion_air_search_response.json")
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	// Parse JSON
	var raw *lionair.Response
	if err := json.Unmarshal(data, &raw); err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	if raw == nil || !raw.Success {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, errors.ErrLionAirNotFound
	}

	for _, f := range raw.Data.AvailableFlights {
		domainFlight, err := f.ToDomainFlightInfo()
		if err != nil {
			logger.ErrorContext(ctx, "Error : ", "err", err)
			return nil, err
		}

		// Filter by input criteria
		if domainFlight.Departure.Airport == input.Origin &&
			domainFlight.Arrival.Airport == input.Destination &&
			domainFlight.Departure.Datetime.Format("2006-01-02") == input.DepartureDate &&
			domainFlight.AvailableSeats >= input.Passengers &&
			strings.EqualFold(domainFlight.CabinClass, input.CabinClass) {
			flights = append(flights, domainFlight)
		}
	}

	return flights, nil
}
