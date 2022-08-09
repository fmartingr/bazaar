package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalModels "github.com/fmartingr/bazaar/internal/models"
	"github.com/fmartingr/bazaar/pkg/manager"
)

type Server struct {
	Http   internalModels.Server
	config *ServerConfig

	cancel context.CancelFunc
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	if s.config.Http.Enabled {
		go func() {
			if err := s.Http.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("error starting server: %s", err)
			}
		}()
	}

	return nil
}

func (s *Server) WaitStop() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signals
	log.Printf("signal %s received, shutting down", sig)

	return s.Stop()
}

func (s *Server) Stop() error {
	s.cancel()

	shuwdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.config.Http.Enabled {
		if err := s.Http.Stop(shuwdownContext); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error shutting down http server: %s", err)
		}
	}

	return nil
}

func NewServer(conf *ServerConfig, m *manager.Manager) *Server {
	server := &Server{
		config: conf,
	}
	if conf.Http.Enabled {
		server.Http = NewHttpServer(conf.Http.Port, m)
	}

	return server
}
