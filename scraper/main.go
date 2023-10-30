package scraper

import (
	"context"
	"fmt"
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
		dp := models.DataPoint{
			Event:    "nicked_server_error",
			Location: "scraper_scheduler_error",
			Details:  err.Error(),
		}
		if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
			log.Println(err)
		}
	}
	s.StartAsync()
}

func trimName(str string) string {
    if (len(str) > 35) {
        return str[0:34] + "..."
    }
    return str
}

func Scrape() {
	ctx := context.Background()

	products, err := models.GetAllProducts(ctx)
	if err != nil {
		dp := models.DataPoint{
			Event:    "nicked_server_error",
			Location: "scraper_get_all_products_error",
			Details:  err.Error(),
		}
		if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
			log.Println(err)
		}
		log.Panic(err)
	}

	for i, product := range products {
		log.Println("Scraping " +
            trimName(product.Name) +
            " " +
            strconv.FormatInt(int64(i+1), 10) +
            "/" +
            strconv.FormatInt(int64(len(products)), 10))

		url, err := url.Parse(product.Url)
		if err != nil {
			dp := models.DataPoint{
				Event:    "nicked_server_error",
				Location: "scraper_parse_url_error",
				Details:  err.Error(),
				Data1:    fmt.Sprintf("url:%s", url),
			}
			if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
				log.Println(err)
			}
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

		if price.Amount == 0 {
			break
		}

		lastPrice := product.Prices[len(product.Prices)-1]
		product.Prices = append(product.Prices, price)

		_, err = models.CreatePrice(price, ctx)
		if err != nil {
			dp := models.DataPoint{
				Event:    "nicked_server_error",
				Location: "scraper_create_price_error",
				Details:  err.Error(),
				Data1:    fmt.Sprintf("productId:%d;price:%v", product.Id, price.Amount),
			}
			if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
				log.Println(err)
			}
			log.Panic(err)
		}

		if !(product.OnSale && price.Amount == lastPrice.Amount) && price.Amount > 0 {
			product.OnSale = price.Amount < lastPrice.Amount

			_, err := models.UpdateProduct(&product, ctx)
			if err != nil {
				dp := models.DataPoint{
					Event:    "nicked_server_error",
					Location: "scraper_amazon_error",
					Details:  err.Error(),
					Data1:    fmt.Sprintf("productId:%d", product.Id),
				}
				if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
					log.Println(err)
				}
				log.Panic(err)
			}

			if product.OnSale {
				for i := 0; i < len(product.Users); i++ {
					user := product.Users[i]
					productSetting, err := models.GetProductSetting(user.Id, product.Id, ctx)
					if err != nil {
						dp := models.DataPoint{
							Event:    "nicked_server_error",
							Location: "scraper_amazon_error",
							Details:  err.Error(),
							Data1:    fmt.Sprintf("userId:%d;productId:%d", user.Id, product.Id),
						}
						if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
							log.Println(err)
						}
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

