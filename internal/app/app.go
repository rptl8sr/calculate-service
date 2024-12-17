package app

type app struct {
}

type App interface {
	Run() error
}

func MustLoad() (App, error) {
	return &app{}, nil
}

func (a *app) Run() error {
	return nil
}
