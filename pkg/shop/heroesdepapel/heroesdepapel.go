package heroesdepapel

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
)

var Domains = []string{"www.heroesdepapel.es"}

type HeroesDePapelShop struct {
	domains []string
}

func (s *HeroesDePapelShop) Get(url string) (*models.Product, error) {
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
	return func() models.Shop {
		shop := HeroesDePapelShop{
			domains: Domains,
		}
		return &shop
	}
}
