package transformers

import (
	"main/core"
	"time"
)

func AddDate(weather *core.Weather) error {
	weather.Timestamp = time.Now().UTC()
	return nil
}
