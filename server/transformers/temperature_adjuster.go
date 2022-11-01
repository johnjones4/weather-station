package transformers

import (
	"main/core"
)

func AdjustTemperature(weather *core.Weather) error {
	weather.TemperatureCalibrationFactor = 1 //TODO
	return nil
}
