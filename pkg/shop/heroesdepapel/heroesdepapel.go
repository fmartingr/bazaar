package heroesdepapel

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
)

var Domains = []string{"www.heroesdepapel.es"}

type HeroesDePapelShop struct {
	models.ShopOptions
	domains []string
}

func (s *HeroesDePapelShop) Get(u *url.URL) (*models.Product, error) {
	body, err := s.ShopOptions.Client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("error during request: %s", err)
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing body: %s", err)
	}

	product := models.Product{
		URL: u.String(),
	}

	doc.Find(".section-product-details").Each(func(i int, s *goquery.Selection) {
		priceText := strings.TrimSpace(s.Find(`.productos-price`).Text())
		priceNum, _ := strconv.ParseFloat(strings.ReplaceAll(strings.Split(priceText, " ")[0], ",", "."), 64)
		// TODO: error logging

		imgURL, _ := s.Find(".carousel-inner .active img").Attr("src")
		// TODO: error logging

		product.Name = strings.TrimSpace(s.Find(".product-title").Nodes[0].FirstChild.Data)
		product.InStock = s.Find(".btn-productos-add-to-cart").Text() != "Agotado"
		product.ImageURL = "https://" + Domains[0] + "/" + imgURL
		product.PriceText = priceText
		product.Price = priceNum
		// releaseDateText := strings.Split(s.Find(".tab-inner-content section-product-description h4").Text(), "A LA VENTA EL")
		// if len(releaseDateText) > 0 {
		// 	releaseDate, _ := time.Parse("2 DE January", releaseDateText[1])
		// 	// TODO: error logging
		// 	product.ReleaseDate = releaseDate
		// }
	})

	return &product, nil
}

func NewHeroesDePapelShopFactory() models.ShopFactory {
	return func(shopOptions models.ShopOptions) models.Shop {
		shop := HeroesDePapelShop{
			ShopOptions: shopOptions,
			domains:     Domains,
		}
		return &shop
	}
}
