package consts

type SortKey string
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"

	SortKeyPrice         SortKey = "price"
	SortKeyDuration      SortKey = "duration"
	SortKeyAirline       SortKey = "airline"
	SortKeyDepartureTime SortKey = "departure_time"
	SortKeyArrivalTime   SortKey = "arrival_time"
	SortKeyBestValue     SortKey = "best_value"
	// SortKeyDeparture     SortKey = "departure"
	// SortKeyArrival       SortKey = "arrival"

)

type FilterKey string

const (
	FilterKeyAirlines        FilterKey = "airlines"
	FilterKeyMinPrice        FilterKey = "min_price"
	FilterKeyMaxPrice        FilterKey = "max_price"
	FilterKeyMinStops        FilterKey = "min_stops"
	FilterKeyMaxStops        FilterKey = "max_stops"
	FilterKeyMinDuration     FilterKey = "min_duration"
	FilterKeyMaxDuration     FilterKey = "max_duration"
	FilterKeyDepartureAfter  FilterKey = "departure_after"
	FilterKeyDepartureBefore FilterKey = "departure_before"
	FilterKeyArrivalAfter    FilterKey = "arrival_after"
	FilterKeyArrivalBefore   FilterKey = "arrival_before"
)
