package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculate-service/internal/controller"
)

func TestCalculate(t *testing.T) {
	testCases := []struct {
		name           string
		expression     string
		expectedResult string
		errorExpected  bool
		expectedCode   int
	}{
		{
			name:           "Valid expression",
			expression:     "2+2",
			expectedResult: "4.000000",
			errorExpected:  false,
			expectedCode:   http.StatusOK,
		},
		{
			name:           "Invalid expression",
			expression:     "2++2",
			expectedResult: "",
			errorExpected:  true,
			expectedCode:   http.StatusUnprocessableEntity,
		},
		{
			name:           "Expression with unsupported operands",
			expression:     "2a+2",
			expectedResult: "",
			errorExpected:  true,
			expectedCode:   http.StatusUnprocessableEntity,
		},
		{
			name:           "Empty expression",
			expression:     "",
			expectedResult: "",
			errorExpected:  true,
			expectedCode:   http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := &CalculatePayload{
				Expression: tc.expression,
			}
			reqBodyBytes, _ := json.Marshal(reqBody)
			req, err := http.NewRequest("POST", "/calculate", bytes.NewReader(reqBodyBytes))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			ctrl := controller.New()
			testHandler := New(ctrl)
			testHandler.Calculate(rec, req)

			if tc.errorExpected {
				if rec.Code != tc.expectedCode {
					t.Fatalf("Expected status %v; got %v", tc.expectedCode, rec.Code)
				}
			} else {
				if rec.Code != http.StatusOK {
					t.Fatalf("expected status 200; got %v", rec.Code)
				}

				var resp CalculateResponse
				err = json.NewDecoder(rec.Body).Decode(&resp)
				if err != nil {
					t.Fatalf("could not decode response: %v", err)
				}

				if resp.Result != tc.expectedResult {
					t.Fatalf("Expected %v, but got %v", tc.expectedResult, resp.Result)
				}
			}
		})
	}
}
