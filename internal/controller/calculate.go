package controller

import (
	"context"
	"errors"

	"calculate-service/pkg/calculator"
)

func (c *controller) Calculate(_ context.Context, expression string) (float64, error) {
	res, err := calculator.Evaluate(expression)

	if err != nil {
		if errors.Is(err, calculator.NewErrUnknown()) {
			return 0, NewServerError(err)
		} else {
			return 0, NewRequestError(err)
		}
	}

	return res, nil
}
