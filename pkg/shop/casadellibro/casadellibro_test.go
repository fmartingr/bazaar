package casadellibro_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/casadellibro"
	"github.com/stretchr/testify/assert"
)

func TestCasaDelLibro_Ok(t *testing.T) {
	shop := casadellibro.NewCasaDelLibroShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl, _ := url.Parse("https://www.casadellibro.com/test/")

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, product.Description)
	assert.Equal(t, "LA DEPENDIENTA", product.Name)
	assert.Equal(t, "https://imagessl0.casadellibro.com/a/l/t5/20/9788416634620.jpg", product.ImageURL)
	assert.Equal(t, 15.96, product.Price)
	assert.Equal(t, "15.96", product.PriceText)
	assert.Equal(t, "2019-01-01T00:00:00Z", product.ReleaseDate.Format(time.RFC3339))
	assert.Equal(t, testUrl.String(), product.URL)
}
