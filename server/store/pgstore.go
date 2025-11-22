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

/*
CREATE TABLE IF NOT EXISTS weather (
  timestamp TIMESTAMP NOT NULL PRIMARY KEY,
  source varchar(128),
  wind_speed REAL NOT NULL,
  vane_direction REAL NOT NULL,
  temperature REAL NOT NULL,
  pressure REAL NOT NULL,
  humidity REAL NOT NULL,
  gas REAL NOT NULL,
  rainfall READ NOT NULL
);
*/

func (s *PGStore) Save(ctx context.Context, w *core.Weather) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO weather (timestamp, source, wind_speed, vane_direction, temperature, pressure, humidity, gas, rainfall) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		w.Timestamp,
		w.Source,
		w.WindSpeed,
		w.VaneDirection,
		w.Temperature,
		w.Pressure,
		w.Humidity,
		w.Gas,
		w.Rainfall,
	)
	return err
}

func (s *PGStore) Get(ctx context.Context, start, end time.Time) ([]core.Weather, error) {
	rows, err := s.pool.Query(ctx, "SELECT timestamp, source, wind_speed, vane_direction, temperature, pressure, humidity, gas, rainfall FROM weather WHERE timestamp >= $1 AND timestamp <= $2 ORDER BY timestamp", start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	weather := make([]core.Weather, 0)
	for rows.Next() {
		var w core.Weather
		err = rows.Scan(
			&w.Timestamp,
			&w.Source,
			&w.WindSpeed,
			&w.VaneDirection,
			&w.Temperature,
			&w.Pressure,
			&w.Humidity,
			&w.Gas,
			&w.Rainfall,
		)
		if err != nil {
			return nil, err
		}
		weather = append(weather, w)
	}
	return weather, nil
}
