package amazon_test

import (
	"net/url"
	"testing"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/amazon"
	"github.com/stretchr/testify/assert"
)

func TestAmazonSpain_Ok(t *testing.T) {
	shop := amazon.NewAmazonShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl, _ := url.Parse("https://www.amazon.es/test/")

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "AGPTEK 50/Pack Combinación de Organizador para Cables (25 Unidades Ajustable Cable Tie Monte Adhesivo + 25 Clips de Cable autoadhesivos), Color Negro", product.Name)
	// assert.Equal(t, "https://www.akiracomics.com/imagenes/poridentidad?identidad=24552a54-365d-4d31-a73e-9fd5f927c3a0&ancho=900&alto=", product.ImageURL)
	assert.Equal(t, 7.99, product.Price)
	assert.Equal(t, "7,99€", product.PriceText)
	assert.Equal(t, testUrl.String(), product.URL)
}
