package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fmartingr/bazaar/pkg/manager"
)

type httpServer struct {
	http    *http.Server
	manager *manager.Manager
}

func (s *httpServer) Start(_ context.Context) error {
	log.Printf("starting http server at %s", s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *httpServer) Stop(ctx context.Context) error {
	log.Println("stoppping http server")
	return s.http.Shutdown(ctx)
}

func NewHttpServer(port int, m *manager.Manager) *httpServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/item", func(rw http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(rw, "ParseForm() err: %v", err)
			return
		}

		product, err := m.Retrieve(r.PostForm.Get("url"))
		if err != nil {
			if errors.Is(err, manager.ErrShopNotFound) {
				rw.WriteHeader(400)
				return
			}

			log.Printf("error for url %s: %s", r.PostForm.Get("url"), err)
			rw.WriteHeader(500)
			return
		}

		payload, _ := json.Marshal(product)

		rw.Header().Add("Content-Type", "application/json")
		if _, err := rw.Write(payload); err != nil {
			log.Printf("error writting response: %s", err)
		}
	})

	mux.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	return &httpServer{
		manager: m,
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}
