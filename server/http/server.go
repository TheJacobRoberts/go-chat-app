package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	router mux.Router

	logger *logrus.Logger

	// Bind address & domain for the server's listener
	// If domain is specified, server is run on TLS using acme/autocert
	Domain, Addr string
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	srv := &Server{
		server: &http.Server{},
		router: *mux.NewRouter(),
		logger: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}

	// TODO: External service for catching panics

	//
	srv.server.Handler = http.HandlerFunc(srv.serveHTTP)

	return srv
}

// UseTLS returns true if the cert and key are specified
func (s *Server) UseTLS() bool {
	return s.Domain != ""
}

// Scheme returns the URL scheme for the server
func (s *Server) Scheme() string {
	if s.UseTLS() {
		return "https"
	}
	return "http"
}

// Port returns the TCP port for the running server
// Useful when testing where we allocate a random port (using ":0")
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server
func (s *Server) URL() string {
	scheme, port := s.Scheme(), s.Port()

	// Use localhost unless a domain is specified
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}

	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}

// Open validates the server options and begins listening on the bind address
func (s *Server) Open() (err error) {
	// Open a listener on our bind address
	if s.Domain != "" {
		s.ln = autocert.NewListener(s.Domain)
	} else {
		if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
			return err
		}
	}

	// Begin serving requests on the listener. We use Server() instead of
	// ListenAndServer() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously
	go s.server.Serve(s.ln)

	return nil
}

// Close gracefully shuts down the server
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("received http request")

	// TODO: Handle forms
	// TODO: Handle custom extensions and set headers and such
	// e.g. https://.../some_file.json => set "Accept", "Content-type", etc.

	s.router.ServeHTTP(w, r)
}

// ListenAndServeTLSRedirect runs an HTTP server on port 80 to redirect users
// to the TLS-enabled port 443 server
func ListenAndServeTLSRedirect(domain string) error {
	return http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+domain, http.StatusFound)
	}))
}
