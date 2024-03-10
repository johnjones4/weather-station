package api

import (
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(store core.Store, transformers []core.Transformer, log *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, log)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", newHealthHandler(store, log))

		r.Route("/weather", func(r chi.Router) {
			r.Post("/", newPostWeatherHandler(store, transformers, log))
			r.Get("/", newGetWeathersHandler(store, log))
		})
	})

	return r
}
