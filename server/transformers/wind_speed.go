package transformers

import (
	"main/core"
)

func CalculateWindSpeed(weather *core.Weather) error {
	weather.WindSpeed = 0 //TODO
	return nil
}
