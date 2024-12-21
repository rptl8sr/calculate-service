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

type ResponseError struct {
	Error string `json:"error"`
}

func (h handler) Calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := CalculatePayload{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		responseError := ResponseError{Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseError)
		return
	}
	defer r.Body.Close()

	if payload.Expression == "" {
		responseError := ResponseError{Error: "'expression' field is required."}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	res, err := h.controller.Calculate(r.Context(), payload.Expression)

	if err != nil {
		var ctrlErr controller.CtrlError
		if errors.As(err, &ctrlErr) {

			switch ctrlErr.Type {
			case controller.ErrRequest:
				responseError := ResponseError{Error: ctrlErr.Error()}
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(responseError)
				return
			case controller.ErrServer:
				responseError := ResponseError{Error: http.StatusText(http.StatusInternalServerError)}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(responseError)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := CalculateResponse{
		Result: fmt.Sprintf("%f", res),
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
