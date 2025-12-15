package domain

import "time"

type AirlineInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type AirportInfo struct {
	Airport   string    `json:"airport"`
	City      string    `json:"city"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

type DurationInfo struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type PriceInfo struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Display  string `json:"display"`
}

type BaggageInfo struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}

type AircraftInfo struct {
	Model string `json:"model"`
	Code  string `json:"code"`
}

type AmenityInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type FlightInfo struct {
	ID             string        `json:"id"`
	Provider       string        `json:"provider"`
	Airline        AirlineInfo   `json:"airline"`
	FlightNumber   string        `json:"flight_number"`
	Departure      AirportInfo   `json:"departure"`
	Arrival        AirportInfo   `json:"arrival"`
	Duration       DurationInfo  `json:"duration"`
	Stops          int           `json:"stops"`
	Price          PriceInfo     `json:"price"`
	AvailableSeats int           `json:"available_seats"`
	CabinClass     string        `json:"cabin_class"`
	Aircraft       *AircraftInfo `json:"aircraft"`
	Amenities      []AmenityInfo `json:"amenities"`
	Baggage        BaggageInfo   `json:"baggage"`
	BestValueScore float64       `json:"best_value_score"`
	// Internal fields not exposed in API
	TotalTripDuration int64 `json:"-"`
}

func (f *FlightInfo) CalculateBestValueScore() {
	f.BestValueScore = float64(f.Price.Amount) / float64(f.Duration.TotalMinutes)
}
