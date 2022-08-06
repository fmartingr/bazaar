package casadellibro_test

import (
	"testing"
	"time"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/casadellibro"
	"github.com/stretchr/testify/assert"
)

func TestCasaDelLibro_Ok(t *testing.T) {
	shop := casadellibro.NewCasaDelLibroShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl := "https://www.casadellibro.com/test/"

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Greater(t, len(product.Description), 100)
	assert.Equal(t, product.Name, "LA DEPENDIENTA")
	assert.Equal(t, product.ImageURL, "https://imagessl0.casadellibro.com/a/l/t5/20/9788416634620.jpg")
	assert.Equal(t, product.Price, 15.96)
	assert.Equal(t, product.PriceText, "15.96")
	assert.Equal(t, product.ReleaseDate.Format(time.RFC3339), "2019-01-01T00:00:00Z")
	assert.Equal(t, product.URL, testUrl)
}
