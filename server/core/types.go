package core

import (
	"context"
	"time"
)

type AnemometerValue struct {
	RotationsPerSecond float64 `json:"rotationsPerSecond"`
	MetersPerSecond    float64 `json:"metersPerSecond"`
}

type Weather struct {
	Timestamp                    time.Time       `json:"timestamp"`
	AnemometerMin                AnemometerValue `json:"anemometerMin"`
	AnemometerMax                AnemometerValue `json:"anemometerMax"`
	AnemometerAverage            AnemometerValue `json:"anemometerAverage"`
	WindSpeed                    float64         `json:"windSpeed"`
	VaneDirection                float64         `json:"vaneDirection"`
	Temperature                  float64         `json:"temperature"`
	TemperatureCalibrationFactor float64         `json:"temperatureCalibrationFactor"`
	Pressure                     float64         `json:"pressure"`
	Gas                          float64         `json:"gas"`
	RelativeHumidity             int             `json:"relativeHumidity"`
}

type Store interface {
	Save(ctx context.Context, w *Weather) error
	Get(ctx context.Context, start, end time.Time) ([]Weather, error)
}

type Transformer func(weather *Weather) error
