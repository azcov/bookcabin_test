package consts

const (
	AirportCGK = "CGK"
	AirportDPS = "DPS"
	AirportSOC = "SOC"
)

var (
	AirportCodeToCity = map[string]string{
		AirportCGK: CityJakarta,
		AirportDPS: CityDenpasar,
		AirportSOC: CitySolo,
	}
)
