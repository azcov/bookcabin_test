package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/stretchr/testify/assert"
)

type MockAirline struct {
	Flights []domain.FlightInfo
	Err     error
}

func (m *MockAirline) SearchFlights(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	time.Sleep(10 * time.Millisecond)
	return m.Flights, nil
}

func TestAirlineProvider_SearchFlights(t *testing.T) {

	flight1 := domain.FlightInfo{ID: "f1", Provider: "airasia"}
	flight2 := domain.FlightInfo{ID: "f2", Provider: "batik"}
	flight3 := domain.FlightInfo{ID: "f3", Provider: "garuda"}

	tests := []struct {
		name            string
		request         domain.SearchRequest
		mocks           map[string]*MockAirline
		expectedTotal   int
		expectedSucc    int
		expectedFail    int
		expectedFlights int
	}{
		{
			name:    "All providers succeed with data",
			request: domain.SearchRequest{Origin: "CGK", Destination: "DPS"},
			mocks: map[string]*MockAirline{
				"airasia": {Flights: []domain.FlightInfo{flight1}, Err: nil},
				"batik":   {Flights: []domain.FlightInfo{flight2}, Err: nil},
				"garuda":  {Flights: []domain.FlightInfo{flight3}, Err: nil},
				"lion":    {Flights: []domain.FlightInfo{}, Err: nil},
			},
			expectedTotal:   3,
			expectedSucc:    4,
			expectedFail:    0,
			expectedFlights: 3,
		},
		{
			name:    "Some providers fail",
			request: domain.SearchRequest{Origin: "CGK", Destination: "DPS"},
			mocks: map[string]*MockAirline{
				"airasia": {Flights: []domain.FlightInfo{flight1}, Err: nil},
				"batik":   {Flights: nil, Err: errors.New("connection timeout")},
				"garuda":  {Flights: []domain.FlightInfo{flight3}, Err: nil},
				"lion":    {Flights: nil, Err: errors.New("internal error")},
			},
			expectedTotal:   2,
			expectedSucc:    2,
			expectedFail:    2,
			expectedFlights: 2,
		},
		{
			name:    "All providers fail",
			request: domain.SearchRequest{Origin: "CGK", Destination: "DPS"},
			mocks: map[string]*MockAirline{
				"airasia": {Flights: nil, Err: errors.New("err")},
				"batik":   {Flights: nil, Err: errors.New("err")},
				"garuda":  {Flights: nil, Err: errors.New("err")},
				"lion":    {Flights: nil, Err: errors.New("err")},
			},
			expectedTotal:   0,
			expectedSucc:    0,
			expectedFail:    4,
			expectedFlights: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ap := &AirlineProvider{
				airAsia:         tt.mocks["airasia"],
				batikAir:        tt.mocks["batik"],
				garudaIndonesia: tt.mocks["garuda"],
				lionAir:         tt.mocks["lion"],
			}

			resp, err := ap.SearchFlights(context.Background(), tt.request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedFlights, len(resp.Flights))
			assert.Equal(t, tt.expectedTotal, resp.Metadata.TotalResults)
			assert.Equal(t, tt.expectedSucc, resp.Metadata.ProvidersSucceeded)
			assert.Equal(t, tt.expectedFail, resp.Metadata.ProvidersFailed)
		})
	}
}
