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
	Scraper.AllowURLRevisit = true

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

	for i, product := range products {
		log.Println("scraping " + product.Name + " " + strconv.FormatInt(int64(i+1), 10) + "/" + strconv.FormatInt(int64(len(products)), 10))

		url, err := url.Parse(product.Url)
		if err != nil {
			log.Panic(err)
		}

		var price models.Price
		price.ProductId = product.Id
		if strings.Contains(url.Host, "amazon") {
			c, p := ScrapeAmazon(product.Url)
			price.Currency = strings.ToUpper(c)
			price.Amount = p
		} else if strings.Contains(url.Host, "wayfair") {
			p := ScrapeWayfair(product.Url)
			price.Amount = p
		}

		log.Println("price.amount", price.Amount)

		if price.Amount == 0 {
            break
		}

		lastPrice := product.Prices[len(product.Prices)-1]
        product.Prices = append(product.Prices, price)

		_, err = models.CreatePrice(price, ctx)
		if err != nil {
			log.Panic(err)
		}

		if !(product.OnSale && price.Amount == lastPrice.Amount) && price.Amount > 0 {
			log.Println("product not on sale with same price, updating product")

			product.OnSale = price.Amount < lastPrice.Amount

			_, err := models.UpdateProduct(&product, ctx)
			if err != nil {
				log.Panic(err)
			}

			if product.OnSale {
				for i := 0; i < len(product.Users); i++ {
					user := product.Users[i]
					productSetting, err := models.GetProductSetting(user.Id, product.Id, ctx)
					if err != nil {
						log.Panic(err)
					}

					if productSetting.Active {
						emailer.SendSaleEmail(user.Email, product)
					}
				}
			}
		}
	}
}

func ScrapeAmazon(url string) (string, float64) {
	var currency string
	var price float64

	Scraper.OnHTML("#tmmSwatches span.a-size-base.a-color-price.a-color-price", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "$") {
			currency = "USD"
		}

		str := strings.ReplaceAll(e.Text, "$", "")
		str = strings.ReplaceAll(str, " ", "")

		log.Println("scraping, found " + str)

		flt, _ := strconv.ParseFloat(str, 64)

		if flt > 0 {
			price = flt
		}
	})

	Scraper.OnHTML("#corePrice_desktop .a-price.a-text-price.apexPriceToPay span", func(e *colly.HTMLElement) {
		if currency == "" && strings.Contains(e.Text, "$") {
			currency = "USD"
		}

		str := strings.ReplaceAll(e.Text, "$", "")
		str = strings.ReplaceAll(str, " ", "")

		log.Println("scraping, found " + str)

		flt, _ := strconv.ParseFloat(str, 64)

		if price <= 0 && flt > 0 {
			price = flt
		}
	})

	Scraper.OnHTML("#corePriceDisplay_desktop_feature_div span.a-price.priceToPay", func(e *colly.HTMLElement) {
		children := e.DOM.Children()

		sym := children.Find("span.a-price-symbol").Text()
		if currency == "" && strings.Contains(sym, "$") {
			currency = "USD"
		}

		str := children.Find("span.a-price-whole").Text()
		str += children.Find("span.a-price-fraction").Text()

		log.Println("scraping, found " + str)

		flt, _ := strconv.ParseFloat(str, 64)

		if price <= 0 && flt > 0 {
			price = flt
		}
	})

	if err := Scraper.Visit(url); err != nil {
		log.Println(err)
	}

	return currency, price
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
