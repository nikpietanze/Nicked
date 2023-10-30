package scraper

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"nicked.io/models"
)

func ScrapeAmazon(url string) (string, float64) {
	var currency string
	var price float64

	Scraper.OnError(func(_ *colly.Response, err error) {
		log.Println(err)
		dp := models.DataPoint{
			Event:    "nicked_server_error",
			Location: "scraper_amazon_error",
			Details:  err.Error(),
			Data1:    url,
		}
        if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
			log.Println(err)
        }
	})

	Scraper.OnHTML("#apex_desktop_newAccordionRow #corePriceDisplay_desktop_feature_div span.a-price.priceToPay", func(e *colly.HTMLElement) {
		children := e.DOM.Children()

		sym := children.Find("span.a-price-symbol").Text()
		if currency == "" && strings.Contains(sym, "$") {
			currency = "USD"
		}

		str := children.Find("span.a-price-whole").Text()
		str += children.Find("span.a-price-fraction").Text()

		flt, _ := strconv.ParseFloat(str, 64)

		if price <= 0 && flt > 0 {
			price = flt
			return
		}
	})

	Scraper.OnHTML("#apex_desktop #corePriceDisplay_desktop_feature_div span.a-price.priceToPay", func(e *colly.HTMLElement) {
		children := e.DOM.Children()

		sym := children.Find("span.a-price-symbol").Text()
		if currency == "" && strings.Contains(sym, "$") {
			currency = "USD"
		}

		str := children.Find("span.a-price-whole").Text()
		str += children.Find("span.a-price-fraction").Text()

		flt, _ := strconv.ParseFloat(str, 64)

		if price <= 0 && flt > 0 {
			price = flt
			return
		}
	})

	Scraper.OnHTML("#tmmSwatches span.a-size-base.a-color-price.a-color-price", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "$") {
			currency = "USD"
		}

		str := strings.ReplaceAll(e.Text, "$", "")
		str = strings.ReplaceAll(str, " ", "")

		flt, _ := strconv.ParseFloat(str, 64)

		if flt > 0 {
			price = flt
			return
		}
	})

	Scraper.OnHTML("#corePrice_desktop .a-price.a-text-price.apexPriceToPay span", func(e *colly.HTMLElement) {
		if currency == "" && strings.Contains(e.Text, "$") {
			currency = "USD"
		}

		str := strings.ReplaceAll(e.Text, "$", "")
		str = strings.ReplaceAll(str, " ", "")

		flt, _ := strconv.ParseFloat(str, 64)

		if price <= 0 && flt > 0 {
			price = flt
			return
		}
	})

    err := Scraper.Visit(url)
    if (err != nil) {
		dp := models.DataPoint{
			Event:    "nicked_server_error",
			Location: "scraper_amazon_visit_error",
			Details:  err.Error(),
			Data1:    url,
		}
        if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
			log.Println(err)
        }
    }

	return currency, price
}
