package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azcov/bookcabin_test/internal/consts"
	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAirlineAggregator
type MockAirlineAggregator struct {
	mock.Mock
}

func (m *MockAirlineAggregator) SearchFlights(ctx context.Context, input domain.SearchRequest) (*domain.SearchResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SearchResponse), args.Error(1)
}

// MockCache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (any, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

func (m *MockCache) Set(key string, value any) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockCache) SetWithExpiration(key string, value any, exp time.Duration) error {
	args := m.Called(key, value, exp)
	return args.Error(0)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func TestFlightService_SerchFlight(t *testing.T) {
	t.Run("CacheHit", func(t *testing.T) {
		mockProvider := new(MockAirlineAggregator)
		mockCache := new(MockCache)
		svc := &flightService{
			airlaneProvider: mockProvider,
			cache:           mockCache,
		}

		req := domain.SearchRequest{
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureDate: "2025-12-25",
			Passengers:    1,
			CabinClass:    "Economy",
		}
		cacheKey := req.ToCacheKey()

		cachedResp := domain.SearchResponse{
			Flights: []domain.FlightInfo{
				{ID: "cached_flight"},
			},
		}

		mockCache.On("Get", cacheKey).Return(cachedResp, nil)

		resp, err := svc.SerchFlight(context.Background(), &req)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Flights))
		assert.Equal(t, "cached_flight", resp.Flights[0].ID)
		assert.True(t, resp.Metadata.CacheHit)

		mockCache.AssertExpectations(t)
		mockProvider.AssertNotCalled(t, "SearchFlights")
	})

	t.Run("CacheMiss_ProviderSuccess", func(t *testing.T) {
		mockProvider := new(MockAirlineAggregator)
		mockCache := new(MockCache)
		svc := &flightService{
			airlaneProvider: mockProvider,
			cache:           mockCache,
		}

		req := domain.SearchRequest{
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureDate: "2025-12-25",
			Passengers:    1,
			CabinClass:    "Economy",
		}
		cacheKey := req.ToCacheKey()

		mockCache.On("Get", cacheKey).Return(nil, errors.New("cache miss"))

		providerResp := &domain.SearchResponse{
			Flights: []domain.FlightInfo{
				{ID: "flight1", Price: domain.PriceInfo{Amount: 1000}, Duration: domain.DurationInfo{TotalMinutes: 60}},
				{ID: "flight2", Price: domain.PriceInfo{Amount: 2000}, Duration: domain.DurationInfo{TotalMinutes: 120}},
			},
		}

		mockProvider.On("SearchFlights", mock.Anything, req).Return(providerResp, nil)
		mockCache.On("Set", cacheKey, mock.Anything).Return(nil)

		resp, err := svc.SerchFlight(context.Background(), &req)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Flights))
		assert.False(t, resp.Metadata.CacheHit)
		assert.Equal(t, 2, resp.Metadata.TotalResults)

		mockCache.AssertExpectations(t)
		mockProvider.AssertExpectations(t)
	})

	t.Run("ProviderError", func(t *testing.T) {
		mockProvider := new(MockAirlineAggregator)
		mockCache := new(MockCache)
		svc := &flightService{
			airlaneProvider: mockProvider,
			cache:           mockCache,
		}

		req := domain.SearchRequest{
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureDate: "2025-12-25",
		}
		cacheKey := req.ToCacheKey()

		mockCache.On("Get", cacheKey).Return(nil, errors.New("miss"))
		mockProvider.On("SearchFlights", mock.Anything, req).Return(nil, errors.New("provider error"))

		resp, err := svc.SerchFlight(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "provider error", err.Error())
	})

	t.Run("WithFilters", func(t *testing.T) {
		mockProvider := new(MockAirlineAggregator)
		mockCache := new(MockCache)
		svc := &flightService{
			airlaneProvider: mockProvider,
			cache:           mockCache,
		}

		req := domain.SearchRequest{
			Filters: []domain.SearchFilter{
				{Key: consts.FilterKeyMaxPrice, Value: 1500.0},
			},
		}

		providerResp := &domain.SearchResponse{
			Flights: []domain.FlightInfo{
				{ID: "cheap", Price: domain.PriceInfo{Amount: 1000}},
				{ID: "expensive", Price: domain.PriceInfo{Amount: 2000}},
			},
		}

		mockCache.On("Get", mock.Anything).Return(nil, errors.New("miss"))
		mockProvider.On("SearchFlights", mock.Anything, req).Return(providerResp, nil)
		mockCache.On("Set", mock.Anything, mock.Anything).Return(nil)

		resp, _ := svc.SerchFlight(context.Background(), &req)

		assert.Equal(t, 1, len(resp.Flights))
		assert.Equal(t, "cheap", resp.Flights[0].ID)
	})

	t.Run("WithSort", func(t *testing.T) {
		mockProvider := new(MockAirlineAggregator)
		mockCache := new(MockCache)
		svc := &flightService{
			airlaneProvider: mockProvider,
			cache:           mockCache,
		}

		req := domain.SearchRequest{
			Sort: domain.SortOption{Key: consts.SortKeyPrice, Order: consts.SortOrderDesc},
		}

		providerResp := &domain.SearchResponse{
			Flights: []domain.FlightInfo{
				{ID: "low", Price: domain.PriceInfo{Amount: 1000}},
				{ID: "high", Price: domain.PriceInfo{Amount: 2000}},
			},
		}

		mockCache.On("Get", mock.Anything).Return(nil, errors.New("miss"))
		mockProvider.On("SearchFlights", mock.Anything, req).Return(providerResp, nil)
		mockCache.On("Set", mock.Anything, mock.Anything).Return(nil)

		resp, _ := svc.SerchFlight(context.Background(), &req)

		assert.Equal(t, 2, len(resp.Flights))
		assert.Equal(t, "high", resp.Flights[0].ID)
		assert.Equal(t, "low", resp.Flights[1].ID)
	})
}
