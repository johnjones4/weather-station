package core

import (
	"context"
	"time"
)

type Weather struct {
	Timestamp                    time.Time `json:"timestamp"`
	AnemometerMin                float64   `json:"anemometerMin"`
	AnemometerMax                float64   `json:"anemometerMax"`
	AnemometerAverage            float64   `json:"anemometerAverage"`
	VaneDirection                float64   `json:"vaneDirection"`
	Temperature                  float64   `json:"temperature"`
	TemperatureCalibrationFactor float64   `json:"temperatureCalibrationFactor"`
	Pressure                     int       `json:"pressure"`
	Gas                          int       `json:"gas"`
	RelativeHumidity             float64   `json:"relativeHumidity"`
}

type WeatherImperial struct {
	Timestamp                    time.Time `json:"timestamp"`
	AnemometerMin                float64   `json:"anemometerMin"`
	AnemometerMax                float64   `json:"anemometerMax"`
	AnemometerAverage            float64   `json:"anemometerAverage"`
	VaneDirection                float64   `json:"vaneDirection"`
	Temperature                  float64   `json:"temperature"`
	TemperatureCalibrationFactor float64   `json:"temperatureCalibrationFactor"`
	Pressure                     float64   `json:"pressure"`
	Gas                          int       `json:"gas"`
	RelativeHumidity             float64   `json:"relativeHumidity"`
}

type Store interface {
	Save(ctx context.Context, w *Weather) error
	Get(ctx context.Context, start, end time.Time) ([]Weather, error)
}

type Transformer func(weather *Weather) error
