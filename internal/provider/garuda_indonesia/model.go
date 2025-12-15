package garudaindonesia

import (
	"time"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/leekchan/accounting"
)

type AirportInfo struct {
	Airport  string    `json:"airport"`
	City     string    `json:"city"`
	Time     time.Time `json:"time"`
	Terminal string    `json:"terminal"`
}

type PriceInfo struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type BaggageInfo struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type TimeInfo struct {
	Airport string    `json:"airport"`
	Time    time.Time `json:"time"`
}

type SegmentInfo struct {
	FlightNumber    string   `json:"flight_number"`
	Departure       TimeInfo `json:"departure"`
	Arrival         TimeInfo `json:"arrival"`
	DurationMinutes int      `json:"duration_minutes"`
	LayoverMinutes  int      `json:"layover_minutes,omitempty"`
}

type FlightInfo struct {
	FlightID        string        `json:"flight_id"`
	Airline         string        `json:"airline"`
	AirlineCode     string        `json:"airline_code"`
	Departure       AirportInfo   `json:"departure"`
	Arrival         AirportInfo   `json:"arrival"`
	DurationMinutes int           `json:"duration_minutes"`
	Stops           int           `json:"stops"`
	Aircraft        string        `json:"aircraft"`
	Price           PriceInfo     `json:"price"`
	AvailableSeats  int           `json:"available_seats"`
	FareClass       string        `json:"fare_class"`
	Baggage         BaggageInfo   `json:"baggage"`
	Amenities       []string      `json:"amenities,omitempty"`
	Segments        []SegmentInfo `json:"segments,omitempty"`
}

func (f *FlightInfo) ToDomainFlightInfo() (domain.FlightInfo, error) {
	// airlineCode, _, _ := util.ParseFlightNumber(f.FlightNumber)
	// departLocation := f.Departure.Time.Location()

	departTs := f.Departure.Time.Unix()
	arriveTs := f.Arrival.Time.Unix()

	formattedDuration := util.FormatDurationMinute(f.DurationMinutes)

	ac := accounting.Accounting{Symbol: f.Price.Currency, Precision: 0}
	formattedPrice := ac.FormatMoney(f.Price.Amount)
	result := domain.FlightInfo{
		ID:       f.FlightID + "_" + f.Airline,
		Provider: f.Airline,
		Airline: domain.AirlineInfo{
			Name: f.Airline,
			Code: f.AirlineCode,
		},
		FlightNumber: f.FlightID,
		Departure: domain.AirportInfo{
			Airport:   f.Departure.Airport,
			City:      f.Departure.City,
			Datetime:  f.Departure.Time,
			Timestamp: departTs,
		},
		Arrival: domain.AirportInfo{
			Airport:   f.Arrival.Airport,
			City:      f.Arrival.City,
			Datetime:  f.Arrival.Time, // convert to depart timezone
			Timestamp: arriveTs,
		},
		Duration: domain.DurationInfo{
			TotalMinutes: f.DurationMinutes,
			Formatted:    formattedDuration,
		},
		Stops: f.Stops,
		Price: domain.PriceInfo{
			Amount:   f.Price.Amount,
			Currency: f.Price.Currency,
			Display:  formattedPrice,
		},
		AvailableSeats: f.AvailableSeats,
		CabinClass:     f.FareClass,
		Aircraft: &domain.AircraftInfo{
			Model: f.Aircraft,
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
	Status  string       `json:"status"`
	Flights []FlightInfo `json:"flights"`
}
