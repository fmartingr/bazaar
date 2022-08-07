package gtmstore

import (
	"fmt"
	"log"
	"net/url"
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

func (s *GTMStoreShop) Get(u *url.URL) (*models.Product, error) {
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
