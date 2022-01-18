package amazon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/utils"
	"github.com/goodsign/monday"
)

var Domains = []string{"www.amazon.es", "www.amazon.com"}

type AmazonShop struct {
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

func (s *AmazonShop) Get(url string) (*models.Product, error) {
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

	var tentativePrice string
	for _, selector := range priceSelectors {
		tentativePrice = strings.TrimSpace(doc.Find(selector).First().Text())
		log.Printf("%s = %s", selector, tentativePrice)
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

	imagesJSON, _ := doc.Find("#main-image-container img").Attr("data-a-dynamic-image")
	// TODO: error handling
	var images map[string]interface{}
	json.Unmarshal([]byte(imagesJSON), &images)
	// TODO: error handling
	var lastImage string
	for key := range images {
		lastImage = key
	}
	product.ImageURL = lastImage

	releaseDateElement := doc.Find(".book_details-publication_date")
	if len(releaseDateElement.Nodes) > 0 {
		releaseDateRaw := releaseDateElement.Parent().Parent().Find(".rpi-attribute-value").Text()

		releaseDate, err := utils.ParseReleaseDate(releaseDateLayoutByDomain[res.Request.URL.Host], strings.TrimSpace(releaseDateRaw), monday.LocaleEsES)
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
	return func() models.Shop {
		shop := AmazonShop{
			domains: Domains,
		}
		return &shop
	}
}
