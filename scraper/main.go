package scraper

import (
	"context"
	"log"
	"net/url"
	"strings"

	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly/v2"

	"nicked.io/emailer"
	"nicked.io/models"
)

var Scraper *colly.Collector

func Init() {
	Scraper = colly.NewCollector()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(10).Minutes().Do(Scrape)
	if err != nil {
		log.Println(err)
	}
    s.StartAsync()
}

func Scrape() {
    ctx := context.Background()

	products, err := models.GetAllProducts(ctx)
	if err != nil {
        log.Panic(err)
	}

	for _, product := range products{
		url, err := url.Parse(product.Url)
		if err != nil {
            log.Panic(err)
		}

        var price models.Price
		price.ProductId = product.Id
		if strings.Contains(url.Host, "amazon") {
			p := ScrapeAmazon(product.Url)
			price.Amount = p
		} else if strings.Contains(url.Host, "wayfair") {
			p := ScrapeWayfair(product.Url)
			price.Amount = p
		}

        lastPrice := product.Prices[len(product.Prices)-1]

        _, err = models.CreatePrice(price, ctx)
        if err != nil {
            log.Panic(err)
        }

        if (product.OnSale && price.Amount == lastPrice.Amount) {
            log.Panic(err)
        } else {
            product.OnSale = price.Amount < lastPrice.Amount;

            _, err := models.UpdateProduct(&product, ctx)
            if (err != nil) {
                log.Panic(err)
            }

            if (product.OnSale) {
                for i := 0; i < len(product.Users); i++ {
                    user := product.Users[i]
                    productSetting, err := models.GetProductSetting(user.Id, product.Id, ctx)
                    if (err != nil) {
                        log.Panic(err)
                    }

                    if (productSetting.Active) {
                        emailer.SendSaleEmail(user.Email, product)
                    }
                }
            }
        }
	}
}

func ScrapeAmazon(url string) float64 {
	var price float64

	Scraper.OnHTML("span.a-size-base.a-color-price.a-color-price", func(e *colly.HTMLElement) {
		str := strings.ReplaceAll(e.Text, "$", "")
		str = strings.ReplaceAll(str, " ", "")

		flt, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Println(err)
		}

		if flt > 0 {
			price = flt
		}
	})

    if err := Scraper.Visit(url); err != nil {
        log.Println(err)
    }
	return price
}

func ScrapeWayfair(url string) float64 {
	var price float64

	Scraper.OnHTML(".a-price-whole", func(e *colly.HTMLElement) {
		flt, err := strconv.ParseFloat(e.Text, 64)
		if err != nil {
			log.Println(err)
		}
		price = flt
	})

	Scraper.OnHTML(".a-price-fraction", func(e *colly.HTMLElement) {
		flt, err := strconv.ParseFloat("0."+e.Text, 64)
		if err != nil {
			log.Println(err)
		}
		price += flt
	})

    if err := Scraper.Visit(url); err != nil {
        log.Println(err)
    }
	return price
}
