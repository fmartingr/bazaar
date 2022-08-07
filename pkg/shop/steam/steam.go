package steam

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
)

var Domains = []string{"store.steampowered.com"}

type SteamShop struct {
	models.ShopOptions
	domains []string
}

func (s *SteamShop) Get(u *url.URL) (*models.Product, error) {
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

	doc.Find(`.page_content_ctn`).Each(func(i int, s *goquery.Selection) {
		priceText := strings.TrimSpace(s.Find(`.game_purchase_action .price`).First().Text())
		priceValue, _ := s.Find(`.game_purchase_price.price`).Attr("data-price-final")
		priceNum, _ := strconv.ParseFloat(strings.ReplaceAll(strings.Split(priceValue, " ")[0], ",", "."), 64)
		// TODO: error logging

		imgURL, _ := s.Find("img.game_header_image_full").Attr("src")
		// TODO: error logging

		product.Name = s.Find("#appHubAppName").Text()
		product.InStock = len(s.Find(".game_area_comingsoon").Nodes) == 0
		product.ImageURL = imgURL
		product.PriceText = priceText
		product.Price = priceNum / 100
		releaseDateText := s.Find(".release_date .date").Text()
		releaseDate, _ := time.Parse("2 Jan, 2006", releaseDateText)
		// TODO: error logging
		product.ReleaseDate = &releaseDate
	})

	return &product, nil
}

func NewSteamShopFactory() models.ShopFactory {
	return func(shopOptions models.ShopOptions) models.Shop {
		shop := SteamShop{
			ShopOptions: shopOptions,
			domains:     Domains,
		}
		return &shop
	}
}
