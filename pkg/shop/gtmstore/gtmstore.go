package gtmstore

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
	"github.com/fmartingr/bazaar/pkg/utils"
)

var Domains = []string{"www.gtm-store.com"}

type GTMStoreShop struct {
	models.ShopOptions
	domains []string
}

func (s *GTMStoreShop) Get(url string) (*models.Product, error) {
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

	doc.Find(`div.primary_block`).Each(func(i int, s *goquery.Selection) {
		priceText := utils.StripLastCharacter(s.Find(".woocommerce-Price-amount.amount").Text())
		priceNum, _ := strconv.ParseFloat(strings.ReplaceAll(priceText, ",", "."), 64)

		var imageURLs []string
		s.Find(".woocommerce-product-gallery__wrapper div").Each(func(i int, s *goquery.Selection) {
			imageURL, exists := s.Find("a").Attr("href")
			if exists {
				log.Println(imageURL)
				imageURLs = append(imageURLs, imageURL)
			}
		})

		product.Name = s.Find(".product_title.entry-title").Text()
		product.InStock = s.Find("p.stock").HasClass("in-stock")
		product.ImageURL = imageURLs[0]
		product.PriceText = priceText
		product.Description = s.Find(".woocommerce-product-details__short-description").Text()
		product.Price = priceNum
	})

	return &product, nil
}

func NewGTMStoreShopFactory() models.ShopFactory {
	return func(shopOptions models.ShopOptions) models.Shop {
		shop := GTMStoreShop{
			ShopOptions: shopOptions,
			domains:     Domains,
		}
		return &shop
	}
}
