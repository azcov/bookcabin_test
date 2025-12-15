package provider

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/errors"
	batikair "github.com/azcov/bookcabin_test/internal/provider/batik_air"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/azcov/bookcabin_test/pkg/ratelimit"
)

type batikAirProvider struct {
	fileDir string
	rl      ratelimit.Limiter
}

func NewBatikAirProvider(fileDir string) AirlineInterface {
	return &batikAirProvider{
		fileDir: fileDir,
		rl:      ratelimit.NewWithDuration(100, time.Second),
	}
}

func (ap *batikAirProvider) SearchFlights(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	if !ap.rl.Allow() {
		return nil, errors.ErrBatikAirRateLimitExceeded
	}
	// Implementation for searching flights from Batik Air
	start := time.Now()
	minDelay := 200
	maxDelay := 400

	defer func() {
		elapsed := time.Since(start)
		delay := util.RandomDuration(minDelay, maxDelay)
		time.Sleep(delay - elapsed)
	}()

	flights, err := ap.callSearch(ctx, input)
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	return flights, nil
}

// callSearch simulates calling the Batik Air search API and returns filtered mock data.
func (ap *batikAirProvider) callSearch(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	var flights []domain.FlightInfo
	// Load mock file
	data, err := os.ReadFile(ap.fileDir + "/batik_air_search_response.json")
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	// Parse JSON
	var raw *batikair.Response
	if err := json.Unmarshal(data, &raw); err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	if raw == nil || raw.Code != 200 {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, errors.ErrBatikAirNotFound
	}

	for _, f := range raw.Results {
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
