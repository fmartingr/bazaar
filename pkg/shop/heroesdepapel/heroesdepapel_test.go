package heroesdepapel_test

import (
	"net/url"
	"testing"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/heroesdepapel"
	"github.com/stretchr/testify/assert"
)

func TesHeroesDePapel_Ok(t *testing.T) {
	shop := heroesdepapel.NewHeroesDePapelShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl, _ := url.Parse("https://www.heroesdepapel.es/test/")

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "Gestalt", product.Name)
	assert.Equal(t, "https://www.heroesdepapel.es/uploads/f.1546.jpg", product.ImageURL)
	assert.Equal(t, 9.95, product.Price)
	assert.Equal(t, "9,95 €", product.PriceText)
	assert.Equal(t, testUrl.String(), product.URL)
}
