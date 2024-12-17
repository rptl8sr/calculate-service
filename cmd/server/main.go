package main

import (
	"fmt"
	"os"

	"calculate-service/internal/app"
)

func main() {
	a, err := app.MustLoad()
	if err != nil {
		fmt.Println(fmt.Errorf("error initializing app: %w", err))
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		defer os.Exit(1)
		fmt.Println("failed to run app", err)
	}
}
