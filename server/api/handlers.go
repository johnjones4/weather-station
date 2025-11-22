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

type weatherResponse struct {
	Readings []core.Weather      `json:"readings"`
	Average  core.WeatherReading `json:"average"`
}

func newGetWeathersHandler(store core.Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		start := readDateParameter(r, "start", now.Add(time.Minute*-30))
		end := readDateParameter(r, "end", now)

		weathers, err := store.Get(r.Context(), start, end)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err, log)
			return
		}

		resp := weatherResponse{
			Readings: make([]core.Weather, 0, len(weathers)),
			Average: core.WeatherReading{
				WindSpeed:     core.Pointer(0.0),
				VaneDirection: core.Pointer(0.0),
				Temperature:   core.Pointer(0.0),
				Pressure:      core.Pointer(0.0),
				Humidity:      core.Pointer(0.0),
				Gas:           core.Pointer(0.0),
				Rainfall:      core.Pointer(0.0),
			},
		}
		counts := core.WeatherReading{
			WindSpeed:     core.Pointer(0.0),
			VaneDirection: core.Pointer(0.0),
			Temperature:   core.Pointer(0.0),
			Pressure:      core.Pointer(0.0),
			Humidity:      core.Pointer(0.0),
			Gas:           core.Pointer(0.0),
			Rainfall:      core.Pointer(0.0),
		}

		units := r.URL.Query().Get("units")
		for _, w := range weathers {
			if units == "imperial" {
				w = transformers.ConvertToImperial(w)
			}
			resp.Readings = append(resp.Readings, w)
			if w.WindSpeed != nil {
				resp.Average.WindSpeed = core.Pointer(*resp.Average.WindSpeed + *w.WindSpeed)
				counts.WindSpeed = core.Pointer(*counts.WindSpeed + 1)
			}
			if w.VaneDirection != nil {
				resp.Average.VaneDirection = core.Pointer(*resp.Average.VaneDirection + *w.VaneDirection)
				counts.VaneDirection = core.Pointer(*counts.VaneDirection + 1)
			}
			if w.Temperature != nil {
				resp.Average.Temperature = core.Pointer(*resp.Average.Temperature + *w.Temperature)
				counts.Temperature = core.Pointer(*counts.Temperature + 1)
			}
			if w.Pressure != nil {
				resp.Average.Pressure = core.Pointer(*resp.Average.Pressure + *w.Pressure)
				counts.Pressure = core.Pointer(*counts.Pressure + 1)
			}
			if w.Humidity != nil {
				resp.Average.Humidity = core.Pointer(*resp.Average.Humidity + *w.Humidity)
				counts.Humidity = core.Pointer(*counts.Humidity + 1)
			}
			if w.Gas != nil {
				resp.Average.Gas = core.Pointer(*resp.Average.Gas + *w.Gas)
				counts.Gas = core.Pointer(*counts.Gas + 1)
			}
			if w.Rainfall != nil {
				resp.Average.Rainfall = core.Pointer(*resp.Average.Rainfall + *w.WindSpeed)
				counts.Rainfall = core.Pointer(*counts.Rainfall + 1)
			}
		}
		resp.Average = core.WeatherReading{
			WindSpeed:     safeDivide(*resp.Average.WindSpeed, *counts.WindSpeed),
			VaneDirection: safeDivide(*resp.Average.VaneDirection, *counts.VaneDirection),
			Temperature:   safeDivide(*resp.Average.Temperature, *counts.Temperature),
			Pressure:      safeDivide(*resp.Average.Pressure, *counts.Pressure),
			Humidity:      safeDivide(*resp.Average.Humidity, *counts.Humidity),
			Gas:           safeDivide(*resp.Average.Gas, *counts.Gas),
			Rainfall:      safeDivide(*resp.Average.Rainfall, *counts.Rainfall),
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
