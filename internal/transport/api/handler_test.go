package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFlightService is a mock implementation of service.FlightInterface
type MockFlightService struct {
	mock.Mock
}

func (m *MockFlightService) SerchFlight(ctx context.Context, input *domain.SearchRequest) (*domain.SearchResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SearchResponse), args.Error(1)
}

func TestHandler_SearchFlights(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockSvc := new(MockFlightService)
		handler := NewHandler(mockSvc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := domain.SearchRequest{
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureDate: "2025-12-25",
			Passengers:    1,
			CabinClass:    "Economy",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/v1/flights/search", bytes.NewBuffer(jsonBytes))

		expectedResp := &domain.SearchResponse{
			SearchCriteria: reqBody,
			Metadata:       domain.SearchMetadata{TotalResults: 1},
			Flights:        []domain.FlightInfo{},
		}

		mockSvc.On("SerchFlight", mock.Anything, &reqBody).Return(expectedResp, nil)

		// Execute
		handler.SearchFlights(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("BadRequest_InvalidJSON", func(t *testing.T) {
		// Setup
		mockSvc := new(MockFlightService)
		handler := NewHandler(mockSvc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/v1/flights/search", bytes.NewBufferString("invalid json"))

		// Execute
		handler.SearchFlights(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "SerchFlight")
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Setup
		mockSvc := new(MockFlightService)
		handler := NewHandler(mockSvc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := domain.SearchRequest{
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureDate: "2025-12-25",
			Passengers:    1,
			CabinClass:    "Economy",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/v1/flights/search", bytes.NewBuffer(jsonBytes))

		expectedErr := errors.New("service failure")
		mockSvc.On("SerchFlight", mock.Anything, &reqBody).Return(nil, expectedErr)

		// Execute
		handler.SearchFlights(c)

		// Assert
		// Assuming httpz.JSONResponse maps errors to 500 or similar
		assert.NotEqual(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
