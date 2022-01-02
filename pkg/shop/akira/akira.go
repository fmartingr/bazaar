package akira

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/gocolly/colly"
)

var Domains = []string{"www.akiracomics.com", "akiracomics.com"}

type AkiraShop struct {
	collector *colly.Collector
	domains   []string
	items     map[string]*models.Product
	itemsMu   sync.RWMutex
}

func (s *AkiraShop) init() {
	s.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	s.collector.OnHTML(`div.panel-ficha-producto div.panel-grupo`, func(h *colly.HTMLElement) {
		priceText := h.ChildText(`[itemprop="price"]`)
		priceNum, err := strconv.ParseFloat(strings.ReplaceAll(strings.Split(priceText, " ")[0], ",", "."), 64)
		if err != nil {
			fmt.Println(err)
		}

		s.itemsMu.Lock()
		s.items[h.Request.URL.String()] = &models.Product{
			Name:      h.ChildText("h1.titulo"),
			InStock:   h.ChildText("span.disponibilidad") == "Disponible",
			URL:       h.Request.URL.String(),
			PriceText: priceText,
			Price:     priceNum,
		}
		s.itemsMu.Unlock()
	})
}

func (s *AkiraShop) Get(url string) (*models.Product, error) {
	if err := s.collector.Visit(url); err != nil {
		return nil, fmt.Errorf("error getting product information: %s", err)
	}
	s.itemsMu.RLock()
	defer s.itemsMu.RUnlock()
	return s.items[url], nil
}

func NewAkiraShopFactory() models.ShopFactory {
	return func(collectorOptions []func(*colly.Collector)) models.Shop {
		shop := AkiraShop{
			collector: colly.NewCollector(collectorOptions...),
			domains:   Domains,
			items:     make(map[string]*models.Product),
			itemsMu:   sync.RWMutex{},
		}
		shop.init()
		return &shop
	}
}
