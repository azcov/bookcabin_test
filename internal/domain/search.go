package domain

import (
	"fmt"
	"strings"

	"github.com/azcov/bookcabin_test/internal/consts"
)

// SearchRequest represents the input parameters for a flight search.
type SearchRequest struct {
	Origin        string         `json:"origin" binding:"required"`
	Destination   string         `json:"destination" binding:"required"`
	DepartureDate string         `json:"departureDate" binding:"required"`
	ReturnDate    *string        `json:"returnDate"`
	Passengers    int            `json:"passengers" binding:"required,gte=1"`
	CabinClass    string         `json:"cabinClass" binding:"required"`
	Filters       []SearchFilter `json:"filters,omitempty"`
	Sort          SortOption     `json:"sort"`
}

type SearchFilter struct {
	Key   consts.FilterKey `json:"key,omitempty"`   // "max_price", "max_stops", "airlines", etc.
	Value any              `json:"value,omitempty"` // value type depends on the filter key
}

func (sr *SearchRequest) ToCacheKey() string {
	key := fmt.Sprintf("search_flight:origin=%s;destination=%s;departureDate=%s;returnDate=%v;passengers=%d;cabinClass=%s;",
		sr.Origin,
		sr.Destination,
		sr.DepartureDate,
		func() string {
			if sr.ReturnDate != nil {
				return *sr.ReturnDate
			}
			return "nil"
		}(),
		sr.Passengers,
		sr.CabinClass,
	)
	var filterKey strings.Builder
	for _, f := range sr.Filters {
		fmt.Fprintf(&filterKey, "%s=%v,", f.Key, f.Value)
	}
	key += filterKey.String() + ";"
	key += fmt.Sprintf("sort_key=%s;sort_order=%s", sr.Sort.Key, sr.Sort.Order)

	return key
}

type SortOption struct {
	Key   consts.SortKey   `json:"key,omitempty"`   // "price", "duration", "departure_time", "arrival_time"
	Order consts.SortOrder `json:"order,omitempty"` // "asc" or "desc"
}

// SearchResponse represents the standardized output of a flight search.
type SearchResponse struct {
	SearchCriteria SearchRequest  `json:"search_criteria"`
	Metadata       SearchMetadata `json:"metadata"`
	Flights        []FlightInfo   `json:"flights"`
}

// type SearchCriteria struct {
// 	Origin        string `json:"origin"`
// 	Destination   string `json:"destination"`
// 	DepartureDate string `json:"departure_date"`
// 	Passengers    int    `json:"passengers"`
// 	CabinClass    string `json:"cabin_class"`
// }

type SearchMetadata struct {
	TotalResults       int  `json:"total_results"`
	ProvidersQueried   int  `json:"providers_queried"`
	ProvidersSucceeded int  `json:"providers_succeeded"`
	ProvidersFailed    int  `json:"providers_failed"`
	SearchTimeMs       int  `json:"search_time_ms"`
	CacheHit           bool `json:"cache_hit"`
}
