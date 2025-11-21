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
			var weather core.Weather
			err := readJson(r, &weather)
			if err != nil {
				errorResponse(w, http.StatusBadRequest, err, log)
				return
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

		resp := map[string]any{}

		units := r.URL.Query().Get("units")
		if units == "imperial" {
			imperialWeathers := make([]core.WeatherImperial, len(weathers))
			for i, w := range weathers {
				imperialWeathers[i] = transformers.ConvertToImperial(w)
			}
			resp["items"] = imperialWeathers
		} else {
			resp["items"] = weathers
		}

		jsonResponse(w, http.StatusOK, resp, log)
	}
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
