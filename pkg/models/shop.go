package models

import "github.com/fmartingr/bazaar/pkg/clients"

type ShopFactory func(baseShop ShopOptions) Shop

type Shop interface {
	Get(url string) (*Product, error)
}

type ShopOptions struct {
	Client clients.Client
}

func NewShopOptions(client clients.Client) ShopOptions {
	return ShopOptions{
		Client: client,
	}
}
