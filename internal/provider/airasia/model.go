package airasia

import (
	"strings"
	"time"

	"github.com/azcov/bookcabin_test/internal/consts"
	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/util"
	"github.com/leekchan/accounting"
)

const ProviderName = "AirAsia"

type StopInfo struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}

type FlightInfo struct {
	FlightCode    string     `json:"flight_code"`
	Airline       string     `json:"airline"`
	FromAirport   string     `json:"from_airport"`
	ToAirport     string     `json:"to_airport"`
	DepartTime    time.Time  `json:"depart_time"`
	ArriveTime    time.Time  `json:"arrive_time"`
	DurationHours float64    `json:"duration_hours"`
	DirectFlight  bool       `json:"direct_flight"`
	PriceIdr      int        `json:"price_idr"`
	Seats         int        `json:"seats"`
	CabinClass    string     `json:"cabin_class"`
	BaggageNote   string     `json:"baggage_note"`
	Stops         []StopInfo `json:"stops,omitempty"`
}

func (f *FlightInfo) ToDomainFlightInfo() (domain.FlightInfo, error) {
	airlineCode, _, _ := util.ParseFlightNumber(f.FlightCode)

	departTs := f.DepartTime.Unix()
	arriveTs := f.ArriveTime.Unix()

	totalMinutes := int(f.DurationHours * 60)
	formattedDuration := util.FormatDurationMinute(totalMinutes)
	ac := accounting.Accounting{Symbol: "IDR", Precision: 0, Format: "%s %v", Thousand: ".", Decimal: ","}
	formattedPrice := ac.FormatMoney(f.PriceIdr)

	baggageInfo := strings.Split(f.BaggageNote, ",")

	result := domain.FlightInfo{
		ID:       f.FlightCode + "_" + f.Airline,
		Provider: f.Airline,
		Airline: domain.AirlineInfo{
			Name: f.Airline,
			Code: airlineCode,
		},
		FlightNumber: f.FlightCode,
		Departure: domain.AirportInfo{
			Airport:   f.FromAirport,
			City:      consts.AirportCodeToCity[f.FromAirport],
			Datetime:  f.DepartTime,
			Timestamp: departTs,
		},
		Arrival: domain.AirportInfo{
			Airport:   f.ToAirport,
			City:      consts.AirportCodeToCity[f.ToAirport],
			Datetime:  f.ArriveTime,
			Timestamp: arriveTs,
		},
		Duration: domain.DurationInfo{
			TotalMinutes: totalMinutes,
			Formatted:    formattedDuration,
		},
		Stops: len(f.Stops),
		Price: domain.PriceInfo{
			Amount:   f.PriceIdr,
			Currency: "IDR",
			Display:  formattedPrice,
		},
		AvailableSeats: f.Seats,
		CabinClass:     f.CabinClass,
		Aircraft:       nil,
		Amenities:      []domain.AmenityInfo{},
		// Baggage: domain.BaggageInfo{
		// 	CarryOn: "Cabin baggage only",
		// 	Checked: "Additional fee",
		// },
	}

	if len(baggageInfo) >= 2 {
		result.Baggage = domain.BaggageInfo{
			CarryOn: baggageInfo[0],
			Checked: baggageInfo[1],
		}
	}
	result.CalculateBestValueScore()
	return result, nil
}

type Response struct {
	Status  string       `json:"status"`
	Flights []FlightInfo `json:"flights"`
}
