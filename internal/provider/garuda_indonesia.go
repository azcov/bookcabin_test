package provider

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/errors"
	garudaindonesia "github.com/azcov/bookcabin_test/internal/provider/garuda_indonesia"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/azcov/bookcabin_test/pkg/ratelimit"
)

type garudaIndonesiaProvider struct {
	fileDir string
	rl      ratelimit.Limiter
}

func NewGarudaIndonesiaProvider(fileDir string) AirlineInterface {
	return &garudaIndonesiaProvider{
		fileDir: fileDir,
		rl:      ratelimit.NewWithDuration(100, time.Second),
	}
}

func (ap *garudaIndonesiaProvider) SearchFlights(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	if !ap.rl.Allow() {
		return nil, errors.ErrGarudaIndonesiaRateLimitExceeded
	}
	// Implementation for searching flights from Garuda Indonesia
	start := time.Now()
	minDelay := 50
	maxDelay := 100

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

// callSearch simulates calling the Garuda Indonesia search API and returns filtered mock data.
func (ap *garudaIndonesiaProvider) callSearch(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	var flights []domain.FlightInfo
	// Load mock file
	data, err := os.ReadFile(ap.fileDir + "/garuda_indonesia_search_response.json")
	if err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	// Parse JSON
	var raw *garudaindonesia.Response
	if err := json.Unmarshal(data, &raw); err != nil {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, err
	}

	if raw == nil || raw.Status != "success" {
		logger.ErrorContext(ctx, "Error : ", "err", err)
		return nil, errors.ErrGarudaIndonesiaNotFound
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
