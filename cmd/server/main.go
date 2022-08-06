package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fmartingr/bazaar/pkg/manager"
	"github.com/fmartingr/bazaar/pkg/shop/akiracomics"
	"github.com/fmartingr/bazaar/pkg/shop/amazon"
	"github.com/fmartingr/bazaar/pkg/shop/gtmstore"
	"github.com/fmartingr/bazaar/pkg/shop/heroesdepapel"
	"github.com/fmartingr/bazaar/pkg/shop/steam"
)

func main() {
	m := manager.NewManager()
	m.Register(akiracomics.Domains, akiracomics.NewAkiraShopFactory())
	m.Register(steam.Domains, steam.NewSteamShopFactory())
	m.Register(heroesdepapel.Domains, heroesdepapel.NewHeroesDePapelShopFactory())
	m.Register(amazon.Domains, amazon.NewAmazonShopFactory())
	m.Register(gtmstore.Domains, gtmstore.NewGTMStoreShopFactory())

	http.HandleFunc("/item", func(rw http.ResponseWriter, r *http.Request) {
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
		rw.Write(payload)
	})

	log.Println("starting server")

	if err := http.ListenAndServe(":5001", http.DefaultServeMux); err != nil {
		log.Printf("Error: %s", err)
	}
}
