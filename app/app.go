package app

import (
	"context"
	"log"
	"os"

	httpserver "go-chat-app/server/http"

	"github.com/sirupsen/logrus"
)

type App struct {
	HTTPServer *httpserver.Server

	Logger *logrus.Logger
}

func NewApp() *App {
	return &App{
		HTTPServer: httpserver.NewServer(),
		Logger: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}
}

func (a *App) Run(ctx context.Context) error {
	if err := a.HTTPServer.Open(); err != nil {
		return err
	}

	if a.HTTPServer.UseTLS() {
		go func() {
			log.Fatal(httpserver.ListenAndServeTLSRedirect(a.HTTPServer.Domain))
		}()
	}

	log.Printf("running: url=%q debug=http://localhost:6060", a.HTTPServer.URL())

	return nil
}

// Close gracefully stops the application
func (a *App) Close() error {
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Close(); err != nil {
			return err
		}
	}
	return nil
}
