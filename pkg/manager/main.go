package manager

import (
	"fmt"
	"net/url"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
)

var ErrShopNotFound = fmt.Errorf("shop not found for domain")

type Manager struct {
	domains map[string]models.Shop
}

func (m *Manager) Register(domains []string, shopFactory models.ShopFactory) {
	baseShop := models.NewShopOptions(clients.NewBasicHttpClient())
	shop := shopFactory(baseShop)

	for _, domain := range domains {
		if _, exists := m.domains[domain]; !exists {
			m.domains[domain] = shop
		}
	}
}

func (m *Manager) GetShop(host string) models.Shop {
	shop, exists := m.domains[host]
	if !exists {
		return nil
	}
	return shop
}

func (m *Manager) Retrieve(productURL string) (*models.Product, error) {
	itemUrl, err := url.Parse(productURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %s", err)
	}

	shop := m.GetShop(itemUrl.Host)
	if shop == nil {
		return nil, ErrShopNotFound
	}

	return shop.Get(itemUrl)
}

func NewManager() Manager {
	return Manager{
		domains: make(map[string]models.Shop),
	}
}
