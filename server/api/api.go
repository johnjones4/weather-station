package api

import (
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(store core.Store, transformers []core.Transformer) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", indexHandler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", newHealthHandler(store))

		r.Route("/weather", func(r chi.Router) {
			r.Post("/", newPostWeatherHandler(store, transformers))
			r.Get("/", newGetWeathersHandler(store))
		})
	})

	return r
}
