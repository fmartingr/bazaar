package steam_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/steam"
	"github.com/stretchr/testify/assert"
)

func TestSteam_Ok(t *testing.T) {
	shop := steam.NewSteamShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl, _ := url.Parse("https://store.steampowered.com/test/")

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "Please, Don’t Touch Anything", product.Name)
	assert.Equal(t, "https://cdn.akamai.steamstatic.com/steam/apps/354240/header.jpg?t=1579122180", product.ImageURL)
	assert.Equal(t, 4.99, product.Price)
	assert.Equal(t, "4,99€", product.PriceText)
	assert.Equal(t, "2015-03-26T00:00:00Z", product.ReleaseDate.Format(time.RFC3339))
	assert.Equal(t, testUrl.String(), product.URL)
}
