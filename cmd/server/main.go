package main

import (
	"github.com/fmartingr/bazaar/internal/server"
	"github.com/fmartingr/bazaar/pkg/manager"
	"github.com/fmartingr/bazaar/pkg/shop/akiracomics"
	"github.com/fmartingr/bazaar/pkg/shop/amazon"
	"github.com/fmartingr/bazaar/pkg/shop/casadellibro"
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
	m.Register(casadellibro.Domains, casadellibro.NewCasaDelLibroShopFactory())

	server := server.NewServer(server.ServerConf{
		HttpPort: 8080,
	}, &m)

	server.Start()

	server.WaitStop()
}
