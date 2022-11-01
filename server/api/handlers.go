package api

import (
	"errors"
	"fmt"
	"main/core"
	"net/http"
	"time"
)

func newPostWeatherHandler(store core.Store, transformers []core.Transformer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-type")
		if ct == "application/json" {
			var weather core.Weather
			err := readJson(r, &weather)
			if err != nil {
				errorResponse(w, http.StatusBadRequest, err)
				return
			}

			for _, t := range transformers {
				err = t(&weather)
				if err != nil {
					errorResponse(w, http.StatusInternalServerError, err)
					return
				}
			}

			err = store.Save(r.Context(), &weather)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err)
				return
			}

			jsonResponse(w, http.StatusOK, weather)

			return
		}
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("unrecognized content type: %s", ct))
	}
}

func newGetWeathersHandler(store core.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := readDateParameter(r, "start")
		end := readDateParameter(r, "end")

		weathers, err := store.Get(r.Context(), start, end)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]any{
			"items": weathers,
		}
		jsonResponse(w, http.StatusOK, resp)
	}
}

func newHealthHandler(store core.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		then := now.Add(time.Minute * -10)
		weathers, err := store.Get(r.Context(), then, now)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		if len(weathers) == 0 {
			errorResponse(w, http.StatusServiceUnavailable, errors.New("no recent data posted"))
			return
		}

		resp := map[string]any{
			"records": len(weathers),
		}
		jsonResponse(w, http.StatusOK, resp)
	}
}
