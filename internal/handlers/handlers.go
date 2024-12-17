package handlers

import (
	"net/http"

	"calculate-service/internal/controller"
)

type handler struct {
	controller controller.Controller
}

func (h handler) Calculate(w http.ResponseWriter, r *http.Request) {
	_ = h.controller.Calculate(r.Context())
}

type Handler interface {
	Calculate(w http.ResponseWriter, r *http.Request)
}

func New(ctrl controller.Controller) Handler {
	return &handler{
		controller: ctrl,
	}
}
