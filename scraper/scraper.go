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
	println("starting scraper process")

	Scraper = colly.NewCollector()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(2).Minutes().Do(Scrape)
	if err != nil {
		log.Println(err)
	}
}

func Scrape() error {
    ctx := context.Background()
	println("starting automated scraping process")

	products, err := models.GetActiveProducts(ctx)
	if err != nil {
        return err
	}

	for _, product := range products{
		println("scraping product: ", product.Name, product.Url)

		url, err := url.Parse(product.Url)
		if err != nil {
            return err
		}

        var price models.Price
		price.ProductId = product.Id
		if strings.Contains(url.Host, "amazon") {
			println("product type: amazon")
			p := ScrapeAmazon(product.Url)
			println("current product price: %s", p)
			price.Amount = p
		} else if strings.Contains(url.Host, "wayfair") {
			println("product type: wayfair")
			p := ScrapeWayfair(product.Url)
			println("current product price: %s", p)
			price.Amount = p
		}

		println("storing product price in db")
        _, err = models.CreatePrice(price, ctx)
        if err != nil {
            return err
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
