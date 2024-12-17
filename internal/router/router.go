package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"calculate-service/internal/controller"
	"calculate-service/internal/handlers"
)

func New(ctrl controller.Controller, apiVersion string) *chi.Mux {
	h := handlers.New(ctrl)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.URLFormat)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/api", func(r chi.Router) {
		r.Route(apiVersion, func(r chi.Router) {
			r.Post("/calculate", h.Calculate)
		})
	})

	return r
}
