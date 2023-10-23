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

	"Nicked/models"
)

var Scraper *colly.Collector

func Init() {
	return
	Scraper = colly.NewCollector()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(2).Minutes().Do(Scrape)
	if err != nil {
		log.Println(err)
	}
}

func Scrape() error {
    ctx := context.Background()

	products, err := models.GetActiveProducts(ctx)
	if err != nil {
        return err
	}

	for _, product := range products{
		url, err := url.Parse(product.Url)
		if err != nil {
            return err
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

        lastPrice, err := models.GetLatestPriceByProduct(product.Id, ctx)
        if (err != nil) {
            return err
        }

        _, err = models.CreatePrice(price, ctx)
        if err != nil {
            return err
        }

        if (product.OnSale && price.Amount == lastPrice.Amount) {
            return nil
        } else {
            product.OnSale = price.Amount < lastPrice.Amount;
            _, err := models.UpdateProduct(&product, ctx)
            if (err != nil) {
                return err
            }
        }
	}
    return nil
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
