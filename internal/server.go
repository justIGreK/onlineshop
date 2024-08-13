package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	maxHeaderBytes = 1 << 20
)

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error during running server:%w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error during shutting down server:%w", err)
	}
	return nil
}
