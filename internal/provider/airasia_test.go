package provider

import (
	"context"
	"testing"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAirAsiaProvider_SearchFlights(t *testing.T) {
	airAsiaProvider := NewAirAsiaProvider("./mock")
	_, err := airAsiaProvider.SearchFlights(context.Background(), domain.SearchRequest{})
	assert.NoError(t, err)

	ctx := context.Background()
	tests := []struct {
		name            string
		ctx             context.Context
		request         domain.SearchRequest
		expectedFlights int
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
			expectedFlights: 3,
			expectedError:   nil,
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
			expectedFlights: 0,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := airAsiaProvider.SearchFlights(context.Background(), tt.request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			// assert.Equal(t, tt.expectedFlights, resp)
			assert.Equal(t, tt.expectedFlights, len(resp), "Expected %d flights, got %d", tt.expectedFlights, len(resp))
			assert.Equal(t, tt.expectedError, err, "Expected error %v, got %v", tt.expectedError, err)
		})
	}
}
