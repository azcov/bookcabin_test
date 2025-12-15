package provider

import (
	"context"
	"testing"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestLionAirProvider_SearchFlights(t *testing.T) {
	lionAirProvider := NewLionAirProvider("./mock")
	_, err := lionAirProvider.SearchFlights(context.Background(), domain.SearchRequest{})
	assert.NoError(t, err)

	tzJkt, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatal(err)
	}
	tzMakassar, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tests := []struct {
		name            string
		ctx             context.Context
		request         domain.SearchRequest
		expectedFlights []domain.FlightInfo
		expectedError   error
	}{
		{
			name: "Has Result",
			ctx:  ctx,
			request: domain.SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "Economy",
			},
			expectedFlights: []domain.FlightInfo{
				{ID: "JT740_Lion Air", Provider: "Lion Air", Airline: domain.AirlineInfo{Name: "Lion Air", Code: "JT"}, FlightNumber: "JT740", Departure: domain.AirportInfo{Airport: "CGK", City: "Jakarta", Datetime: time.Date(2025, time.December, 15, 5, 30, 0, 0, tzJkt), Timestamp: 1765751400}, Arrival: domain.AirportInfo{Airport: "DPS", City: "Denpasar", Datetime: time.Date(2025, time.December, 15, 8, 15, 0, 0, tzMakassar), Timestamp: 1765757700}, Duration: domain.DurationInfo{TotalMinutes: 105, Formatted: "1h 45m"}, Stops: 0, Price: domain.PriceInfo{Amount: 950000, Currency: "IDR"}, AvailableSeats: 45, CabinClass: "ECONOMY", Aircraft: nil, Amenities: []domain.AmenityInfo{}, Baggage: domain.BaggageInfo{CarryOn: "7 kg", Checked: "20 kg"}, TotalTripDuration: 0},
				{ID: "JT742_Lion Air", Provider: "Lion Air", Airline: domain.AirlineInfo{Name: "Lion Air", Code: "JT"}, FlightNumber: "JT742", Departure: domain.AirportInfo{Airport: "CGK", City: "Jakarta", Datetime: time.Date(2025, time.December, 15, 11, 45, 0, 0, tzJkt), Timestamp: 1765773900}, Arrival: domain.AirportInfo{Airport: "DPS", City: "Denpasar", Datetime: time.Date(2025, time.December, 15, 14, 35, 0, 0, tzMakassar), Timestamp: 1765780500}, Duration: domain.DurationInfo{TotalMinutes: 110, Formatted: "1h 50m"}, Stops: 0, Price: domain.PriceInfo{Amount: 890000, Currency: "IDR"}, AvailableSeats: 38, CabinClass: "ECONOMY", Aircraft: nil, Amenities: []domain.AmenityInfo{}, Baggage: domain.BaggageInfo{CarryOn: "7 kg", Checked: "20 kg"}, TotalTripDuration: 0},
				{ID: "JT650_Lion Air", Provider: "Lion Air", Airline: domain.AirlineInfo{Name: "Lion Air", Code: "JT"}, FlightNumber: "JT650", Departure: domain.AirportInfo{Airport: "CGK", City: "Jakarta", Datetime: time.Date(2025, time.December, 15, 16, 20, 0, 0, tzJkt), Timestamp: 1765790400}, Arrival: domain.AirportInfo{Airport: "DPS", City: "Denpasar", Datetime: time.Date(2025, time.December, 15, 21, 10, 0, 0, tzMakassar), Timestamp: 1765804200}, Duration: domain.DurationInfo{TotalMinutes: 230, Formatted: "3h 50m"}, Stops: 1, Price: domain.PriceInfo{Amount: 780000, Currency: "IDR"}, AvailableSeats: 52, CabinClass: "ECONOMY", Aircraft: nil, Amenities: []domain.AmenityInfo{}, Baggage: domain.BaggageInfo{CarryOn: "7 kg", Checked: "20 kg"}, TotalTripDuration: 0}},
			expectedError: nil,
		},
		{
			name: "No Result",
			ctx:  ctx,
			request: domain.SearchRequest{
				Origin:        "BDO",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "Economy",
			},
			expectedFlights: []domain.FlightInfo{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := lionAirProvider.SearchFlights(context.Background(), tt.request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			// assert.Equal(t, tt.expectedFlights, resp)
			assert.Equal(t, len(tt.expectedFlights), len(resp), "Expected %d flights, got %d", len(tt.expectedFlights), len(resp))
			assert.Equal(t, tt.expectedError, err, "Expected error %v, got %v", tt.expectedError, err)
		})
	}
}
