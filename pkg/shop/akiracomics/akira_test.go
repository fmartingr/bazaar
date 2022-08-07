package akiracomics_test

import (
	"net/url"
	"testing"

	"github.com/fmartingr/bazaar/pkg/clients"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/shop/akiracomics"
	"github.com/stretchr/testify/assert"
)

func TestAkiraComics_Ok(t *testing.T) {
	shop := akiracomics.NewAkiraShopFactory()(models.NewShopOptions(clients.NewMockClient()))

	testUrl, _ := url.Parse("https://www.akiracomics.com/test/")

	product, err := shop.Get(testUrl)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "TODO PARECÍA DEMASIADO FÁCIL.\n\nLos amigos y aliados de Rimuru siguen desconcertados por su ascenso, pero algo los preocupa aún más: el regreso de Veldora, el temido Dragón Tormenta. Al mismo tiempo, la lentitud y la cautela de los súbditos de Rimuru hacen presagiar futuras luchas mayores que están por venir y alteran el estado de poder que se ha creado.", product.Description)
	assert.Equal(t, "AQUELLA VEZ QUE ME CONVERTI EN SLIME VOL.16 [RUSTICA]", product.Name)
	assert.Equal(t, "https://www.akiracomics.com/imagenes/poridentidad?identidad=24552a54-365d-4d31-a73e-9fd5f927c3a0&ancho=900&alto=", product.ImageURL)
	assert.Equal(t, 8.55, product.Price)
	assert.Equal(t, "8,55 €", product.PriceText)
	assert.Equal(t, testUrl.String(), product.URL)
}
