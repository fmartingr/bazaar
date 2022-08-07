package gtmstore_test

import (
	"net/url"
	"testing"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/gtmstore"
	"github.com/stretchr/testify/assert"
)

func TestGtmStore_Ok(t *testing.T) {
	shop := gtmstore.NewGTMStoreShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl, _ := url.Parse("https://www.gtm-store.com/test/")

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(product)

	assert.NotEmpty(t, product.Description)
	assert.Equal(t, "Metroid: Misión Omega", product.Name)
	assert.Equal(t, "https://www.gtm-store.com/wp-content/uploads/2022/07/GTM-Metroid-4.jpeg", product.ImageURL)
	assert.Equal(t, 8.99, product.Price)
	assert.Equal(t, "8,99", product.PriceText)
	assert.Equal(t, testUrl.String(), product.URL)
}
