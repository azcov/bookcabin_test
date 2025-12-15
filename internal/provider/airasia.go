package provider

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/errors"
	"github.com/azcov/bookcabin_test/internal/provider/airasia"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/azcov/bookcabin_test/pkg/ratelimit"
)

type airAsiaProvider struct {
	fileDir string
	rl      ratelimit.Limiter
}

func NewAirAsiaProvider(fileDir string) AirlineInterface {
	return &airAsiaProvider{
		fileDir: fileDir,
		rl:      ratelimit.NewWithDuration(100, time.Second),
	}
}

func (ap *airAsiaProvider) SearchFlights(ctx context.Context, input domain.SearchRequest) (flights []domain.FlightInfo, err error) {
	if !ap.rl.Allow() {
		return nil, errors.ErrAirAsiaRateLimitExceeded
	}
	// Implementation for searching flights from AirAsia
	start := time.Now()
	minDelay := 50
	maxDelay := 150
	failRate := 0.10

	defer func() {
		elapsed := time.Since(start)
		// Simulate delay: 50–150ms

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

	// Simulate failure: 10% fail rate
	if util.RandomFailure(failRate) {
		return nil, errors.ErrAirAsiaInternalError
	}
	flights, err = ap.callSearch(ctx, input)
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	return flights, nil
}

// callSearch simulates calling the AirAsia search API and returns filtered mock data.
func (ap *airAsiaProvider) callSearch(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	var flights []domain.FlightInfo
	// Load mock file
	data, err := os.ReadFile(ap.fileDir + "/airasia_search_response.json")
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	// logger.InfoContext(ctx, "AirAsia Mock Data: ", "data", string(data))
	// Parse JSON
	var raw *airasia.Response
	if err := json.Unmarshal(data, &raw); err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	if raw == nil || raw.Status != "ok" {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, errors.ErrAirAsiaNotFound
	}

	for _, f := range raw.Flights {
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
