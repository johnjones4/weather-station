package api

import (
	_ "embed"
	"errors"
	"fmt"
	"main/core"
	"main/transformers"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

//go:embed index.html
var index []byte

func indexHandler(w http.ResponseWriter, r *http.Request, log *zap.SugaredLogger) {
	w.Header().Set("Content-type", "text/html")
	if r.URL.Query().Get("read") == "" {
		w.Write(index)
		return
	}
	index1, err := os.ReadFile("./api/index.html")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, log)
		return
	}
	w.Write(index1)
}

func newPostWeatherHandler(store core.Store, transformers []core.Transformer, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-type")
		if ct == "application/json" {
			var request core.WeatherPayload
			err := readJson(r, &request)
			if err != nil {
				errorResponse(w, http.StatusBadRequest, err, log)
				return
			}

			weather := core.Weather{
				WeatherPayload: request,
				Timestamp:      time.Now(),
			}

			for _, t := range transformers {
				err = t(&weather)
				if err != nil {
					errorResponse(w, http.StatusInternalServerError, err, log)
					return
				}
			}

			err = store.Save(r.Context(), &weather)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err, log)
				return
			}

			jsonResponse(w, http.StatusOK, weather, log)

			return
		}
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("unrecognized content type: %s", ct), log)
	}
}

type weatherAverage struct {
	Period        string              `json:"period"`
	Start         time.Time           `json:"start"`
	Averages      core.WeatherReading `json:"averages"`
	RainfallTotal float64             `json:"rainfallTotal"`
}

type weatherResponse struct {
	Readings []core.Weather   `json:"readings"`
	Periods  []weatherAverage `json:"periods"`
}

func newGetWeathersHandler(store core.Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		start := readDateParameter(r, "start", now.Add(time.Hour*-12))
		end := readDateParameter(r, "end", now)

		weathers, err := store.Get(r.Context(), start, end)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err, log)
			return
		}

		allPeriods := map[string]time.Time{
			"5min":  end.Add(time.Minute * -5),
			"1hr":   end.Add(-time.Hour),
			"12hrs": end.Add(time.Hour * -12),
		}
		periods := make(map[string][2]core.WeatherReading)
		rainfallTotals := make(map[string]float64)
		zero := core.WeatherReading{
			WindSpeed:     core.Pointer(0.0),
			VaneDirection: core.Pointer(0.0),
			Temperature:   core.Pointer(0.0),
			Pressure:      core.Pointer(0.0),
			Humidity:      core.Pointer(0.0),
			Gas:           core.Pointer(0.0),
			Rainfall:      core.Pointer(0.0),
		}
		for key, start := range allPeriods {
			if start.Before(start) {
				continue
			}
			periods[key] = [2]core.WeatherReading{zero, zero}
			rainfallTotals[key] = 0
		}

		resp := weatherResponse{
			Readings: make([]core.Weather, 0, len(weathers)),
			Periods:  make([]weatherAverage, 0, len(periods)),
		}

		units := r.URL.Query().Get("units")
		for _, w := range weathers {
			if units == "imperial" {
				w = transformers.ConvertToImperial(w)
			}
			resp.Readings = append(resp.Readings, w)
			for key, info := range periods {
				if w.Timestamp.Before(allPeriods[key]) {
					continue
				}
				counts := info[0]
				average := info[1]
				if w.WindSpeed != nil {
					average.WindSpeed = core.Pointer(*average.WindSpeed + *w.WindSpeed)
					counts.WindSpeed = core.Pointer(*counts.WindSpeed + 1)
				}
				if w.VaneDirection != nil {
					average.VaneDirection = core.Pointer(*average.VaneDirection + *w.VaneDirection)
					counts.VaneDirection = core.Pointer(*counts.VaneDirection + 1)
				}
				if w.Temperature != nil {
					average.Temperature = core.Pointer(*average.Temperature + *w.Temperature)
					counts.Temperature = core.Pointer(*counts.Temperature + 1)
				}
				if w.Pressure != nil {
					average.Pressure = core.Pointer(*average.Pressure + *w.Pressure)
					counts.Pressure = core.Pointer(*counts.Pressure + 1)
				}
				if w.Humidity != nil {
					average.Humidity = core.Pointer(*average.Humidity + *w.Humidity)
					counts.Humidity = core.Pointer(*counts.Humidity + 1)
				}
				if w.Gas != nil {
					average.Gas = core.Pointer(*average.Gas + *w.Gas)
					counts.Gas = core.Pointer(*counts.Gas + 1)
				}
				if w.Rainfall != nil {
					average.Rainfall = core.Pointer(*average.Rainfall + *w.Rainfall)
					counts.Rainfall = core.Pointer(*counts.Rainfall + 1)
					rainfallTotals[key] += *average.Rainfall
				}
				periods[key] = [2]core.WeatherReading{counts, average}
			}
		}
		for key, info := range periods {
			resp.Periods = append(resp.Periods, weatherAverage{
				Period: key,
				Start:  allPeriods[key],
				Averages: core.WeatherReading{
					WindSpeed:     safeDivide(*info[1].WindSpeed, *info[0].WindSpeed),
					VaneDirection: safeDivide(*info[1].VaneDirection, *info[0].VaneDirection),
					Temperature:   safeDivide(*info[1].Temperature, *info[0].Temperature),
					Pressure:      safeDivide(*info[1].Pressure, *info[0].Pressure),
					Humidity:      safeDivide(*info[1].Humidity, *info[0].Humidity),
					Gas:           safeDivide(*info[1].Gas, *info[0].Gas),
					Rainfall:      safeDivide(*info[1].Rainfall, *info[0].Rainfall),
				},
				RainfallTotal: rainfallTotals[key],
			})
		}

		jsonResponse(w, http.StatusOK, resp, log)
	}
}

func safeDivide(a float64, b float64) *float64 {
	if b == 0 {
		return nil
	}
	return core.Pointer(a / b)
}

func newHealthHandler(store core.Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		then := now.Add(time.Minute * -10)
		weathers, err := store.Get(r.Context(), then, now)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err, log)
			return
		}

		if len(weathers) == 0 {
			errorResponse(w, http.StatusServiceUnavailable, errors.New("no recent data posted"), log)
			return
		}

		resp := map[string]any{
			"records": len(weathers),
		}
		jsonResponse(w, http.StatusOK, resp, log)
	}
}
