package casadellibro

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fmartingr/bazaar/pkg/models"
)

const (
	releaseDatePropertyName        = "Fecha de lanzamiento"
	releaseDatePropertyValueFormat = "02/01/2006"
)

var Domains = []string{"www.casadellibro.com"}

type CasaDelLibroShop struct {
	domains     []string
	priceRegexp *regexp.Regexp
}

func (s *CasaDelLibroShop) Get(url string) (*models.Product, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error retrieving url: %s", err)
	}

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

	// Price
	seoJSON := doc.Find(`[vmid="seo"][data-body="true"]`).Text()

	priceSubmatch := s.priceRegexp.FindSubmatch([]byte(seoJSON))
	if len(priceSubmatch) != 2 {
		return nil, fmt.Errorf("error parsing price, expected one match, found %d", len(priceSubmatch)/2)
	}

	priceNum, _ := strconv.ParseFloat(strings.ReplaceAll(string(priceSubmatch[1]), ",", "."), 64)
	product.PriceText = string(priceSubmatch[1])
	product.Price = priceNum

	// Release date (& meta)
	rightBlocks := doc.Find(".border-left.col-md-4.col-12")
	if rightBlocks.Length() >= 2 {
		rightBlocks.Eq(1).Find(".row").Each(func(i int, s *goquery.Selection) {
			cols := s.Find(".col")
			propertyNamme := strings.Trim(cols.First().Text(), ":")
			if propertyNamme == releaseDatePropertyName {
				propertyValue := cols.Eq(1).Text()

				releaseDate, err := time.Parse(releaseDatePropertyValueFormat, propertyValue)
				if err != nil {
					log.Printf("error parsing release date: %s", err)
					return
				}

				product.ReleaseDate = &releaseDate
			}
		})
	}

	doc.Find(`main.v-main`).Each(func(i int, s *goquery.Selection) {
		imageURL, _ := s.Find(".product-image").Attr("src")
		product.Name = strings.TrimSpace(s.Find("h1.text-h5.mb-4").Text())
		product.InStock = s.Find("p.stock").HasClass("in-stock")
		product.ImageURL = imageURL
		product.Description = strings.TrimSpace(s.Find(".resume-body").Text())
	})

	return &product, nil
}

func NewCasaDelLibroShopFactory() models.ShopFactory {
	return func() models.Shop {
		r, err := regexp.Compile(`Price\"\:\"([\d+\.]+)`)
		if err != nil {
			log.Println(err)
		}

		shop := CasaDelLibroShop{
			domains:     Domains,
			priceRegexp: r,
		}
		return &shop
	}
}
