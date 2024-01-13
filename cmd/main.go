package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go-chat-app/app"
)

func main() {
	// Setup signal handlers
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	app := app.NewApp()

	if err := app.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
