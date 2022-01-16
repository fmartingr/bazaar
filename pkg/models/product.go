package models

import "time"

type Product struct {
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	ImageURL    string     `json:"image_url"`
	InStock     bool       `json:"in_stock"`
	PriceText   string     `json:"price_text"`
	Price       float64    `json:"price"`
	ReleaseDate *time.Time `json:"release_date"`
}
