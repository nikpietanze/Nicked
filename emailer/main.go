package emailer

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"nicked.io/config"
	"nicked.io/models"
)

func SendSaleEmail(recipient string, product *models.Product) {
	// TODO: currency symbol based on product.Currency
	price := strconv.FormatFloat(product.Prices[len(product.Prices)-1].Amount, 'f', -1, 64)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	html, err := os.ReadFile(dir + "/emailer/templates/price-drop.html")
	if err != nil {
		panic(err)
	}

	body := string(html)
    //body = strings.ReplaceAll(body, "{{logo_url}}", "https://nicked.io/static/img/logo.png")
    body = strings.ReplaceAll(body, "{{logo_url}}", "https://ebhsrec.stripocdn.email/content/guids/CABINET_6d4f2213bba54903a0d6ab253b3f66b9077afefa2804034a2e515cf9737a4f9e/images/group_36.png")
    body = strings.ReplaceAll(body, "{{product_img_url}}", product.ImageUrl)
	body = strings.ReplaceAll(body, "{{product_name}}", product.Name)
	body = strings.ReplaceAll(body, "{{product_price}}", "$"+price)
    body = strings.ReplaceAll(body, "{{product_url}}", product.Url)

	subject := fmt.Sprintf("Price Drop Alert: %s", product.Name)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		mime + "\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", config.EMAIL_ADDRESS, config.EMAIL_PASSWORD, config.EMAIL_HOST_ADDRESS)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", config.EMAIL_HOST_ADDRESS, config.EMAIL_HOST_PORT),
		auth,
		config.EMAIL_ADDRESS,
		[]string{recipient},
		msg)
	if err != nil {
		log.Panic(err)
	}
	log.Println("email is successfully sent to " + recipient)
}
