package transformers

import "main/core"

func ConvertToImperial(weather core.Weather) core.Weather {
	w := core.Weather{
		Timestamp: weather.Timestamp,
		WeatherPayload: core.WeatherPayload{
			WeatherReading: core.WeatherReading{
				VaneDirection: weather.VaneDirection,
				Gas:           weather.Gas,
				Humidity:      weather.Humidity,
			},
		},
	}
	if weather.WindSpeed != nil {
		s := mpsToMph(*weather.WindSpeed)
		w.WindSpeed = &s
	}
	if weather.Temperature != nil {
		t := (*weather.Temperature)*1.8 + 32
		w.Temperature = &t
	}
	if weather.Pressure != nil {
		p := (*weather.Pressure) / 100 * 0.02953
		w.Pressure = &p
	}
	return w
}

func mpsToMph(in float64) float64 {
	return in * 2.236936
}
