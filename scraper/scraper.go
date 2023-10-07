package scraper

import (
	"log"
	"net/url"
	"strings"

	"reflect"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly/v2"
	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

var Scraper *colly.Collector

func Init() {
	Scraper = colly.NewCollector()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(24).Hours().Do(Scrape)
	if err != nil {
		log.Println(err)
	}
}

func Scrape() {
	ctxPtr := reflect.New(reflect.TypeOf(new(iris.Context)))
	ctx := ctxPtr.Elem().Interface().(iris.Context)

	items, err := models.GetActiveItems(ctx)
	if err != nil {
		log.Println(err)
	}

	for _, item := range items {
		url, err := url.Parse(item.Url)
		if err != nil {
		    log.Println(err)
		}

		price := new(models.Price)
		if strings.Contains(url.Host, "amazon") {
			p := ScrapeAmazon(item.Url)
			price.Amount = p
		} else if strings.Contains(url.Host, "wayfair") {
			p := ScrapeWayfair(item.Url)
			price.Amount = p
		}

		models.CreatePrice(price, ctx)
	}
}

func ScrapeAmazon(url string) float64 {
	var price float64

	Scraper.OnHTML(".SFPrice span", func(e *colly.HTMLElement) {
		flt, err := strconv.ParseFloat(e.Text, 64)
		if err != nil {
		    log.Println(err)
		}
		price = flt
	})

	Scraper.Visit(url)
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

	Scraper.Visit(url)
	return price
}
