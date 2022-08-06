package akiracomics

import (
	"fmt"
	"net/http"
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

func (s *AkiraShop) Get(url string) (*models.Product, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error retrieving url: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error retrieving url: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing body: %s", err)
	}

	product := models.Product{
		URL: url,
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
