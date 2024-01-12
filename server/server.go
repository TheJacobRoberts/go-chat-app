package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
	logger *logrus.Logger
}

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Test is what we usually do"))
}

func NewServer(handler http.Handler, port int, logger *logrus.Logger) *Server {
	r := mux.NewRouter()
	r.HandleFunc("/test", TestEndpoint).Methods(http.MethodGet)

	srv := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: r,
		},
		logger: logger,
	}

	return srv
}

func (s *Server) GracefulShutdown(killtime time.Duration) {
	stopCh := make(chan os.Signal, 1)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), killtime)
	defer cancel()

	s.logger.Info("shutting down handlers...")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}

// TODO: Read https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
func (s *Server) Start() error {
	s.logger.Infof("listening server on %v\n", s.server.Addr)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
