package batikair

import (
	"time"

	"github.com/azcov/bookcabin_test/internal/consts"
	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/leekchan/accounting"
)

const ProviderName = "Batik Air"

var (
	classToCabinClass = map[string]string{
		"Y": "Economy",
		// "X":   "Business",
		// "Z": "First Class",
	}

	TimeFormat = "2006-01-02T15:04:05Z0700" // "RFC3339     = "2006-01-02T15:04:05Z07:00""
)

type FareInfo struct {
	BasePrice    int    `json:"basePrice"`
	Taxes        int    `json:"taxes"`
	TotalPrice   int    `json:"totalPrice"`
	CurrencyCode string `json:"currencyCode"`
	Class        string `json:"class"`
}

type ConnectionInfo struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}

type FlightInfo struct {
	FlightNumber      string           `json:"flightNumber"`
	AirlineName       string           `json:"airlineName"`
	AirlineIATA       string           `json:"airlineIATA"`
	Origin            string           `json:"origin"`
	Destination       string           `json:"destination"`
	DepartureDateTime string           `json:"departureDateTime"`
	ArrivalDateTime   string           `json:"arrivalDateTime"`
	TravelTime        string           `json:"travelTime"`
	NumberOfStops     int              `json:"numberOfStops"`
	Fare              FareInfo         `json:"fare"`
	SeatsAvailable    int              `json:"seatsAvailable"`
	AircraftModel     string           `json:"aircraftModel"`
	BaggageInfo       string           `json:"baggageInfo"`
	OnboardServices   []string         `json:"onboardServices"`
	Connections       []ConnectionInfo `json:"connections,omitempty"`
}

func (f *FlightInfo) ToDomainFlightInfo() (domain.FlightInfo, error) {
	// airlineCode, _, _ := util.ParseFlightNumber(f.FlightNumber)
	departTime, err := time.Parse(TimeFormat, f.DepartureDateTime)
	if err != nil {
		logger.Error("Error : ", "err", err)
		return domain.FlightInfo{}, err
	}
	arriveTime, err := time.Parse(TimeFormat, f.ArrivalDateTime)
	if err != nil {
		logger.Error("Error : ", "err", err)
		return domain.FlightInfo{}, err
	}
	departTs := departTime.Unix()
	arriveTs := arriveTime.Unix()

	// totalMinutes := int(f.TravelTime * 60)
	// formattedDuration := util.FormatDurationMinute(totalMinutes)

	totalMinutes := arriveTime.Sub(departTime).Minutes()
	ac := accounting.Accounting{Symbol: f.Fare.CurrencyCode, Precision: 0}
	formattedPrice := ac.FormatMoney(f.Fare.TotalPrice)

	result := domain.FlightInfo{
		ID:       f.FlightNumber + "_" + f.AirlineName,
		Provider: f.AirlineName,
		Airline: domain.AirlineInfo{
			Name: f.AirlineName,
			Code: f.AirlineIATA,
		},
		FlightNumber: f.FlightNumber,
		Departure: domain.AirportInfo{
			Airport:   f.Origin,
			City:      consts.AirportCodeToCity[f.Origin],
			Datetime:  departTime,
			Timestamp: departTs,
		},
		Arrival: domain.AirportInfo{
			Airport:   f.Destination,
			City:      consts.AirportCodeToCity[f.Destination],
			Datetime:  arriveTime,
			Timestamp: arriveTs,
		},
		Duration: domain.DurationInfo{
			TotalMinutes: int(totalMinutes),
			Formatted:    f.TravelTime,
		},
		Stops: f.NumberOfStops,
		Price: domain.PriceInfo{
			Amount:   f.Fare.TotalPrice,
			Currency: f.Fare.CurrencyCode,
			Display:  formattedPrice,
		},
		AvailableSeats: f.SeatsAvailable,
		CabinClass:     classToCabinClass[f.Fare.Class],
		Aircraft: &domain.AircraftInfo{
			Model: f.AircraftModel,
			Code:  "",
		},
		Amenities: []domain.AmenityInfo{},
		Baggage: domain.BaggageInfo{
			CarryOn: "Cabin baggage only",
			Checked: "Additional fee",
		},
	}

	return result, nil
}

type Response struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Results []FlightInfo `json:"results"`
}
