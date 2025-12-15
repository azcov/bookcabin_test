package lionair

import (
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/leekchan/accounting"
)

var (
	TimeFormat = "2006-01-02T15:04:05"
)

type Carrier struct {
	Name string `json:"name"`
	Iata string `json:"iata"`
}

type Airport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Route struct {
	From Airport `json:"from"`
	To   Airport `json:"to"`
}

type Schedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type Pricing struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
	FareType string `json:"fare_type"`
}

type BaggageAllowance struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}

type Services struct {
	WifiAvailable    bool             `json:"wifi_available"`
	MealsIncluded    bool             `json:"meals_included"`
	BaggageAllowance BaggageAllowance `json:"baggage_allowance"`
}

type LayoverInfo struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}

type FlightInfo struct {
	ID         string        `json:"id"`
	Carrier    Carrier       `json:"carrier"`
	Route      Route         `json:"route"`
	Schedule   Schedule      `json:"schedule"`
	FlightTime int           `json:"flight_time"`
	IsDirect   bool          `json:"is_direct"`
	Pricing    Pricing       `json:"pricing"`
	SeatsLeft  int           `json:"seats_left"`
	PlaneType  string        `json:"plane_type"`
	Services   Services      `json:"services"`
	StopCount  int           `json:"stop_count,omitempty"`
	Layovers   []LayoverInfo `json:"layovers,omitempty"`
}

func (f *FlightInfo) ToDomainFlightInfo() (domain.FlightInfo, error) {
	// Convert times
	departTz, err := time.LoadLocation(f.Schedule.DepartureTimezone)
	if err != nil {
		logger.Error("Error : ", "err", err)
		return domain.FlightInfo{}, err
	}
	departTime, err := time.ParseInLocation(TimeFormat, f.Schedule.Departure, departTz)
	if err != nil {
		logger.Error("Error : ", "err", err)
		return domain.FlightInfo{}, err
	}

	arriveTz, err := time.LoadLocation(f.Schedule.ArrivalTimezone)
	if err != nil {
		logger.Error("Error : ", "err", err)
		return domain.FlightInfo{}, err
	}
	arriveTime, err := time.ParseInLocation(TimeFormat, f.Schedule.Arrival, arriveTz)
	if err != nil {
		logger.Error("Error : ", "err", err)
		return domain.FlightInfo{}, err
	}

	departTs := departTime.Unix()
	arriveTs := arriveTime.Unix()

	// Format duration
	formattedDuration := util.FormatDurationMinute(f.FlightTime)

	ac := accounting.Accounting{Symbol: f.Pricing.Currency, Precision: 0, Format: "%s %v", Thousand: ".", Decimal: ","}
	formattedPrice := ac.FormatMoney(f.Pricing.Total)

	result := domain.FlightInfo{
		ID:       f.ID + "_" + f.Carrier.Name,
		Provider: f.Carrier.Name,

		Airline: domain.AirlineInfo{
			Name: f.Carrier.Name,
			Code: f.Carrier.Iata,
		},

		FlightNumber: f.ID,

		Departure: domain.AirportInfo{
			Airport:   f.Route.From.Code,
			City:      f.Route.From.City,
			Datetime:  departTime,
			Timestamp: departTs,
		},

		Arrival: domain.AirportInfo{
			Airport:   f.Route.To.Code,
			City:      f.Route.To.City,
			Datetime:  arriveTime,
			Timestamp: arriveTs,
		},

		Duration: domain.DurationInfo{
			TotalMinutes: f.FlightTime,
			Formatted:    formattedDuration,
		},

		Stops: f.StopCount,

		Price: domain.PriceInfo{
			Amount:   f.Pricing.Total,
			Currency: f.Pricing.Currency,
			Display:  formattedPrice,
		},

		AvailableSeats: f.SeatsLeft,
		CabinClass:     f.Pricing.FareType,

		Aircraft: &domain.AircraftInfo{
			Model: f.PlaneType,
			Code:  "",
		},

		Amenities: []domain.AmenityInfo{},

		Baggage: domain.BaggageInfo{
			CarryOn: f.Services.BaggageAllowance.Cabin,
			Checked: f.Services.BaggageAllowance.Hold,
		},
	}

	if f.Services.WifiAvailable {
		result.Amenities = append(result.Amenities, domain.AmenityInfo{
			Type:        "WiFi",
			Description: "WiFi Availabe",
		})
	}
	if f.Services.MealsIncluded {
		result.Amenities = append(result.Amenities, domain.AmenityInfo{
			Type:        "Meal",
			Description: "Meal Included",
		})
	}
	result.CalculateBestValueScore()
	return result, nil
}

type Response struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []FlightInfo `json:"available_flights"`
	} `json:"data"`
}
