package amazon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/utils"
	"github.com/goodsign/monday"
)

var Domains = []string{"www.amazon.es", "www.amazon.com"}

type AmazonShop struct {
	models.ShopOptions
	domains []string
}

var priceSelectors = []string{
	"#buybox span.a-color-price",
	"#tp_price_block_total_price_ww",
	"#price",
	".priceToPay .a-offscreen",
}

var releaseDateLayoutByDomain = map[string]string{
	Domains[0]: "2 January 2006",
	Domains[1]: "January 2, 2006",
}

func (s *AmazonShop) Get(u *url.URL) (*models.Product, error) {
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

	var tentativePrice string
	for _, selector := range priceSelectors {
		tentativePrice = strings.TrimSpace(doc.Find(selector).First().Text())
		if tentativePrice != "" {
			priceNum, err := strconv.ParseFloat(utils.ExtractPrice(tentativePrice), 64)
			if err != nil {
				log.Println(err)
			} else {
				product.PriceText = tentativePrice
				product.Price = priceNum
				break
			}
		}
	}

	product.Name = strings.TrimSpace(doc.Find("#productTitle").Text())

	imagesJSON, exists := doc.Find("#main-image-container img").Attr("data-a-dynamic-image")
	if !exists {
		log.Printf("Can't find image for %s", u.String())
	}
	var images map[string]interface{}
	if err := json.Unmarshal([]byte(imagesJSON), &images); err != nil {
		log.Printf("error unmarshalling: %s", err)
	}
	var lastImage string
	for key := range images {
		lastImage = key
	}
	product.ImageURL = lastImage

	releaseDateElement := doc.Find(".book_details-publication_date")
	if len(releaseDateElement.Nodes) > 0 {
		releaseDateRaw := releaseDateElement.Parent().Parent().Find(".rpi-attribute-value").Text()

		releaseDate, err := utils.ParseReleaseDate(releaseDateLayoutByDomain[u.Host], strings.TrimSpace(releaseDateRaw), monday.LocaleEsES)
		if err != nil {
			log.Println(err)
		} else {
			product.ReleaseDate = releaseDate
		}
	}

	product.InStock = !strings.Contains(doc.Find("#availability").Text(), "No disponible")

	return &product, nil
}

func NewAmazonShopFactory() models.ShopFactory {
	return func(shopOptions models.ShopOptions) models.Shop {
		shop := AmazonShop{
			ShopOptions: shopOptions,
			domains:     Domains,
		}
		return &shop
	}
}
