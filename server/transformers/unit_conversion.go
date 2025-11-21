package transformers

import "main/core"

func ConvertToImperial(weather core.Weather) core.WeatherImperial {
	w := core.WeatherImperial{
		Timestamp:                    weather.Timestamp,
		VaneDirection:                weather.VaneDirection,
		TemperatureCalibrationFactor: weather.TemperatureCalibrationFactor,
		Gas:                          weather.Gas,
		RelativeHumidity:             weather.RelativeHumidity,
	}
	w.AnemometerAverage = mpsToMph(weather.AnemometerAverage)
	w.AnemometerMax = mpsToMph(weather.AnemometerMax)
	w.AnemometerMin = mpsToMph(weather.AnemometerMin)
	w.Temperature = weather.Temperature*1.8 + 32
	w.Pressure = float64(weather.Pressure) / 100 * 0.02953
	return w
}

func mpsToMph(in float64) float64 {
	return in * 2.236936
}
