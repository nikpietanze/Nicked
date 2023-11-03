package scraper

import (
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

func ScrapeWayfair(url string) float64 {
    c := colly.NewCollector()
	var price float64

	c.OnHTML(".a-price-whole", func(e *colly.HTMLElement) {
		flt, err := strconv.ParseFloat(e.Text, 64)
		if err != nil {
			log.Println(err)
		}
		price = flt
	})

	c.OnHTML(".a-price-fraction", func(e *colly.HTMLElement) {
		flt, err := strconv.ParseFloat("0."+e.Text, 64)
		if err != nil {
			log.Println(err)
		}
		price += flt
	})

	if err := c.Visit(url); err != nil {
		log.Println(err)
	}
	return price
}
