package handlers

import (
	"net/http"

	"calculate-service/internal/controller"
)

type handler struct {
	controller controller.Controller
}

type Handler interface {
	Calculate(w http.ResponseWriter, r *http.Request)
}

func New(ctrl controller.Controller) Handler {
	return &handler{
		controller: ctrl,
	}
}
