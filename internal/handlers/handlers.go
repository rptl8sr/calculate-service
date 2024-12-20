package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"calculate-service/internal/controller"
)

type CalculatePayload struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result string `json:"result"`
}

type handler struct {
	controller controller.Controller
}

func (h handler) Calculate(w http.ResponseWriter, r *http.Request) {
	expression := CalculatePayload{}

	err := json.NewDecoder(r.Body).Decode(&expression)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.controller.Calculate(r.Context(), expression.Expression)
	if err != nil {
		var ctrlErr controller.CtrlError
		if errors.As(err, &ctrlErr) {
			switch ctrlErr.Type {
			case controller.ErrRequest:
				http.Error(w, ctrlErr.Error(), http.StatusUnprocessableEntity)
			case controller.ErrServer:
				http.Error(w, ctrlErr.Error(), http.StatusInternalServerError)
			}
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	response := CalculateResponse{
		Result: fmt.Sprintf("%f", res),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type Handler interface {
	Calculate(w http.ResponseWriter, r *http.Request)
}

func New(ctrl controller.Controller) Handler {
	return &handler{
		controller: ctrl,
	}
}
