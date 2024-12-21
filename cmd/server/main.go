package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"calculate-service/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		cancel()
	}()

	a, err := app.MustLoad()
	if err != nil {
		fmt.Println(fmt.Errorf("error initializing app: %w", err))
		os.Exit(1)
	}

	if err = a.Run(ctx); err != nil {
		fmt.Println("failed to run app", err)
		os.Exit(1)
	}
}
