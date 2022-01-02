package manager

import (
	"fmt"

	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/gocolly/colly"
)

type Manager struct {
	domains          map[string]models.Shop
	collectorOptions []func(*colly.Collector)
}

func (m *Manager) Register(domains []string, shopFactory models.ShopFactory) error {
	options := m.collectorOptions
	options = append(options, colly.AllowedDomains(domains...))
	shop := shopFactory(options)

	for _, domain := range domains {
		if _, exists := m.domains[domain]; exists {
			return fmt.Errorf("domain %s is already registered", domain)
		} else {
			m.domains[domain] = shop
		}
	}

	return nil
}

func (m *Manager) Get(host string) models.Shop {
	shop, exists := m.domains[host]
	if !exists {
		return nil
	}

	return shop
}

func NewManager() Manager {
	return Manager{
		collectorOptions: []func(*colly.Collector){
			colly.UserAgent("bazaar/0.0.1"),
		},
		domains: make(map[string]models.Shop),
	}
}
