package emailer

import (
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"nicked.io/config"
	"nicked.io/models"
)

func SendSaleEmail(recipient string, product models.Product) {
	// TODO: currency symbol based on product.Currency
	price := strconv.FormatFloat(product.Prices[len(product.Prices)-1].Amount, 'f', -1, 64)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	html, err := os.Open(dir + "/emailer/templates/price-drop.html")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := html.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := html.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
	}

	body := string(buf)
    body = strings.ReplaceAll(body, "{{product_img_url}}", product.ImageUrl)
	body = strings.ReplaceAll(body, "{{product_name}}", product.Name)
	body = strings.ReplaceAll(body, "{{product_price}}", "$"+price)

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
