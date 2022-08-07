package akiracomics

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
)

var Domains = []string{"www.akiracomics.com", "akiracomics.com"}

type AkiraShop struct {
	models.ShopOptions
	domains []string
}

func (s *AkiraShop) Get(u *url.URL) (*models.Product, error) {
	body, err := s.ShopOptions.Client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("error during request: %s", err)
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing body: %s", err)
	}

	description, err := doc.Find(`[itemprop="description"]`).Html()
	if err != nil {
		// TODO error logging
	} else {
		description = strings.Replace(description, "<br/>", "\n", -1)
	}

	product := models.Product{
		URL:         u.String(),
		Description: description,
	}

	doc.Find(`div.panel-ficha-producto div.panel-grupo`).Each(func(i int, s *goquery.Selection) {
		priceText := s.Find(`[itemprop="price"]`).Text()
		priceNum, _ := strconv.ParseFloat(strings.ReplaceAll(strings.Split(priceText, " ")[0], ",", "."), 64)
		// TODO: error logging

		// Javascript injects a img.zoomImg without the height/widht paramenters, we could remove the parameters
		// from the URL we get. It's most likely that the "thumbnail" is enough for most use cases.
		imgURL, _ := s.Find("a.portada img").Attr("src")
		// TODO: error logging

		product.Name = s.Find("h1.titulo").Text()
		product.InStock = s.Find("span.disponibilidad").Text() == "Disponible"
		product.ImageURL = "https://" + Domains[0] + imgURL
		product.PriceText = priceText
		product.Price = priceNum
	})

	doc.Find(`.panel-descripcion-propiedades`).Each(func(i int, s *goquery.Selection) {
		releaseDateText := s.Find(".fechaedicion .valor-propiedad").Text()
		releaseDate, _ := time.Parse("02/01/2006", releaseDateText)
		// TODO: error logging
		product.ReleaseDate = &releaseDate
	})
	return &product, nil
}

func NewAkiraShopFactory() models.ShopFactory {
	return func(shopOptions models.ShopOptions) models.Shop {
		shop := AkiraShop{
			ShopOptions: shopOptions,
			domains:     Domains,
		}
		return &shop
	}
}
