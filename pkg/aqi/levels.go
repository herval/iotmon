package aqi

type AirQuality string

const (
	Good                     AirQuality = "green"
	Moderate                            = "yellow"
	UnhealthySensitiveGroups            = "orange"
	Unhealthy                           = "red"
	VeryUnhealthy                       = "purple"
	Hazardous                           = "maroon"
)

// Mapping levels https://www.airnow.gov/aqi/aqi-basics/
func AqiForLevel(level int) AirQuality {
	if level <= 50 {
		return Good
	} else if level >= 51 && level <= 100 {
		return Moderate
	} else if level >= 101 && level <= 150 {
		return UnhealthySensitiveGroups
	} else if level >= 151 && level <= 200 {
		return Unhealthy
	} else if level >= 201 && level <= 300 {
		return VeryUnhealthy
	}
	return Hazardous
}
