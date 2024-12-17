package controller

import (
	"context"
)

type controller struct{}

func (c *controller) Calculate(context context.Context) error {
	return nil
}

type Controller interface {
	Calculate(context context.Context) error
}

func New() Controller {
	return &controller{}
}
