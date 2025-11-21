package core

import (
	"context"
	"time"
)

type WeatherReading struct {
	WindSpeed     *float64 `json:"windSpeed"`
	VaneDirection *float64 `json:"vaneDirection"`
	Temperature   *float64 `json:"temperature"`
	Pressure      *float64 `json:"pressure"`
	Humidity      *float64 `json:"humidity"`
	Gas           *float64 `json:"gas"`
	Rainfall      *float64 `json:"rainfall"`
}

type WeatherPayload struct {
	Source string `json:"source"`
	WeatherReading
}

type Weather struct {
	Timestamp time.Time `json:"timestamp"`
	WeatherPayload
}

type Store interface {
	Save(ctx context.Context, w *Weather) error
	Get(ctx context.Context, start, end time.Time) ([]Weather, error)
}

type Transformer func(weather *Weather) error
