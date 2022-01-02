package models

import "github.com/gocolly/colly"

type ShopFactory func(collectorOptions []func(*colly.Collector)) Shop

type Shop interface {
	Get(url string) (*Product, error)
}
