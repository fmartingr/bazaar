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

type ServerConf struct {
	HttpPort int
}

type Server struct {
	http internalModels.Server

	cancel context.CancelFunc
}

func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	go func() {
		if err := s.http.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting server: %s", err)
		}
	}()

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

	if err := s.http.Stop(shuwdownContext); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error shutting down http server: %s", err)
	}

	return nil
}

func NewServer(serverConf ServerConf, m *manager.Manager) *Server {
	return &Server{
		http: NewHttpServer(serverConf.HttpPort, m),
	}
}
