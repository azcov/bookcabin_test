package provider

import (
	"context"

	"github.com/azcov/bookcabin_test/internal/domain"
)

type AirlineInterface interface {
	SearchFlights(ctx context.Context, input domain.SearchRequest) ([]domain.FlightInfo, error)
}

type AirlineAggregator interface {
	SearchFlights(ctx context.Context, input domain.SearchRequest) (*domain.SearchResponse, error)
}
