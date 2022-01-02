package models

type Product struct {
	Name      string  `json:"name"`
	URL       string  `json:"url"`
	InStock   bool    `json:"in_stock"`
	PriceText string  `json:"price_text"`
	Price     float64 `json:"price"`
}
