package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/fmartingr/bazaar/pkg/manager"
	"github.com/fmartingr/bazaar/pkg/shop/akira"
)

func main() {
	m := manager.NewManager()
	m.Register(akira.Domains, akira.NewAkiraShopFactory())

	http.HandleFunc("/item", func(rw http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(rw, "ParseForm() err: %v", err)
			return
		}

		itemUrl, err := url.Parse(r.PostForm.Get("url"))
		if err != nil {
			rw.WriteHeader(400)
			return
		}

		shop := m.Get(itemUrl.Host)
		if shop == nil {
			rw.WriteHeader(400)
			return
		}
		product, err := shop.Get(itemUrl.String())
		if err != nil {
			rw.WriteHeader(500)
			return
		}

		payload, _ := json.Marshal(product)

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(payload)
	})

	if err := http.ListenAndServe(":5001", http.DefaultServeMux); err != nil {
		log.Printf("Error: %s", err)
	}
}
