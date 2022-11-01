package store

import (
	"context"
	"main/core"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(ctx context.Context, conStr string) (core.Store, error) {
	pool, err := pgxpool.Connect(ctx, conStr)
	if err != nil {
		return nil, err
	}

	return &PGStore{pool}, nil
}

func (s *PGStore) Save(ctx context.Context, w *core.Weather) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO weather (tstamp, anemometer_rotations, wind_speed, vein_direction, temperature, temperature_calibration_factor, gas, relative_humidity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		w.Timestamp,
		w.AnemometerRotations,
		w.WindSpeed,
		w.VeinDirection,
		w.Temperature,
		w.TemperatureCalibrationFactor,
		w.Gas,
		w.RelativeHumidity,
	)
	return err
}

func (s *PGStore) Get(ctx context.Context, start, end time.Time) ([]core.Weather, error) {
	rows, err := s.pool.Query(ctx, "SELECT tstamp, anemometer_rotations, wind_speed, vein_direction, temperature, temperature_calibration_factor, gas, relative_humidity FROM weather WHERE tstamp >= $1 AND tstamp <= $2 ORDER BY tstamp", start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	weather := make([]core.Weather, 0)
	for rows.Next() {
		var w core.Weather
		err = rows.Scan(
			&w.Timestamp,
			&w.AnemometerRotations,
			&w.WindSpeed,
			&w.VeinDirection,
			&w.Temperature,
			&w.TemperatureCalibrationFactor,
			&w.Gas,
			&w.RelativeHumidity,
		)
		if err != nil {
			return nil, err
		}
		weather = append(weather, w)
	}
	return weather, nil
}
