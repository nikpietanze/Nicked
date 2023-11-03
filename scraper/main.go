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

func Init() {
	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()
	_, err := s.Every(1).Minutes().Do(Scrape)
	if err != nil {
		dp := models.DataPoint{
			Event:    "nicked_server_error",
			Location: "scraper_scheduler_error",
			Details:  err.Error(),
		}
		if err := models.CreateDataPoint(&dp, context.Background()); err != nil {
			log.Println(err)
		}
		log.Fatalln("error scheduling job", err)
	}
	s.StartAsync()
}

func Scrape() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowURLRevisit(),
        colly.MaxBodySize(1024 * 1024 * 1.5),
	)
    if err := c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 16}); err != nil {
		log.Println(err)
	}
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
    c.DisableCookies()

	ctx := context.Background()

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL)
	})

	c.OnHTML("#centerCol", func(e *colly.HTMLElement) {
		cont := true

		e.ForEachWithBreak("#apex_desktop_newAccordionRow #corePriceDisplay_desktop_feature_div span.a-price.priceToPay", func(_ int, h *colly.HTMLElement) bool {
			children := h.DOM.Children()

			str := children.Find("span.a-price-whole").Text()
			str += children.Find("span.a-price-fraction").Text()
			flt, _ := strconv.ParseFloat(str, 64)

			if flt > 0 {
				processPrice(flt, e.Request.URL.Query().Get("npid"), ctx)
				cont = false
				return false
			}
			return true
		})

		if cont {
			e.ForEachWithBreak("#apex_desktop #corePriceDisplay_desktop_feature_div span.a-price.priceToPay", func(_ int, h *colly.HTMLElement) bool {
				children := h.DOM.Children()

				str := children.Find("span.a-price-whole").Text()
				str += children.Find("span.a-price-fraction").Text()
				flt, _ := strconv.ParseFloat(str, 64)

				if flt > 0 {
					processPrice(flt, e.Request.URL.Query().Get("npid"), ctx)
					cont = false
					return false
				}
				return true
			})
		}

		if cont {
			e.ForEachWithBreak("#tmmSwatches span.a-size-base.a-color-price.a-color-price", func(_ int, e *colly.HTMLElement) bool {
				str := strings.ReplaceAll(e.Text, "$", "")
				str = strings.ReplaceAll(str, " ", "")
				flt, _ := strconv.ParseFloat(str, 64)

				if flt > 0 {
					processPrice(flt, e.Request.URL.Query().Get("npid"), ctx)
					cont = false
					return false
				}
				return true
			})
		}

		if cont {
			e.ForEachWithBreak("#corePrice_desktop .a-price.a-text-price.apexPriceToPay span", func(_ int, e *colly.HTMLElement) bool {
				str := strings.ReplaceAll(e.Text, "$", "")
				str = strings.ReplaceAll(str, " ", "")
				flt, _ := strconv.ParseFloat(str, 64)

				if flt > 0 {
					processPrice(flt, e.Request.URL.Query().Get("npid"), ctx)
					cont = false
					return false
				}
				return true
			})
		}

		if cont {
			e.ForEachWithBreak("#corePrice_desktop #snsPriceRow #snsDetailPagePrice #sns-base-price", func(_ int, e *colly.HTMLElement) bool {
				str := strings.ReplaceAll(e.Text, "$", "")
				str = strings.ReplaceAll(str, " ", "")
				flt, _ := strconv.ParseFloat(str, 64)

				if flt > 0 {
					processPrice(flt, e.Request.URL.Query().Get("npid"), ctx)
					return true
				}
				return false
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println(err)
		dp := models.DataPoint{
			Event:    "nicked_server_error",
			Location: "scraper_amazon_error",
			Details:  err.Error(),
			Data1:    r.Request.URL.String(),
		}
		if err := models.CreateDataPoint(&dp, ctx); err != nil {
			log.Println(err)
		}
	})

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

	for i := 0; i < len(products); i++ {
		product := products[i]

		url, err := url.Parse(product.Url)
		if err != nil {
			log.Println(err)
		}
		if err := c.Visit(fmt.Sprintf("%s?npid=%d",
			url.Scheme+"://"+url.Host+url.Path,
			products[i].Id,
		)); err != nil {
			log.Println(err)
		}
	}

	c.Wait()
}

func processPrice(p float64, npid string, ctx context.Context) {
	log.Println("productId:", npid, "price:", p)
	if p == 0 {
		return
	}

	productId, err := strconv.ParseInt(npid, 10, 64)
	if err != nil {
		log.Println(err)
	}

	product, err := models.GetProduct(productId, ctx)
	if err != nil {
		log.Println(err)
	}

	if product.Id == productId {
		price := models.Price{
			Amount:    p,
			Currency:  "USD",
			ProductId: product.Id,
		}
		lastPrice := product.Prices[len(product.Prices)-1]
		product.Prices = append(product.Prices, price)

		_, err := models.CreatePrice(price, ctx)
		if err != nil {
			dp := models.DataPoint{
				Event:    "nicked_server_error",
				Location: "scraper_create_price_error",
				Details:  err.Error(),
				Data1:    fmt.Sprintf("productId:%d;price:%v", product.Id, price.Amount),
			}
			if err := models.CreateDataPoint(&dp, ctx); err != nil {
				log.Println(err)
			}
			log.Panic(err)
		}

		if !(product.OnSale && price.Amount == lastPrice.Amount) && price.Amount > 0 {
			product.OnSale = price.Amount < lastPrice.Amount

			_, err := models.UpdateProduct(product, ctx)
			if err != nil {
				dp := models.DataPoint{
					Event:    "nicked_server_error",
					Location: "scraper_amazon_error",
					Details:  err.Error(),
					Data1:    fmt.Sprintf("productId:%d", product.Id),
				}
				if err := models.CreateDataPoint(&dp, ctx); err != nil {
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
						if err := models.CreateDataPoint(&dp, ctx); err != nil {
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
