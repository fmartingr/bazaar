package models

type ShopFactory func() Shop

type Shop interface {
	Get(url string) (*Product, error)
}
