package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"nicked.io/models"
)

type ProductJSON struct {
	Email    string
	Name     string
	ImageUrl string
	OnSale   bool
	Price    []PriceJSON
	Sku      string
	Store    string
	Url      string
	UserId   int64
}

type PriceJSON struct {
	Id       string
	Amount   string
	Currency string
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product id")
		// TODO: Send DP
	}

	productId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// Send DP
	}

	product, err := models.GetProduct(productId, c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	return c.JSON(http.StatusOK, product)
}

func CreateProduct(c echo.Context) error {
	var productJSON ProductJSON
	if err := c.Bind(&productJSON); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product")
		// TODO: Send DP
	}

	caser := cases.Title(language.English)

	product, err := models.CreateProduct(
		&models.Product{
			Name:     caser.String(productJSON.Name),
			ImageUrl: productJSON.ImageUrl,
			OnSale:   productJSON.OnSale || false,
			Sku:      productJSON.Sku,
			Store:    strings.ToLower(productJSON.Store),
			Url:      strings.ToLower(productJSON.Url),
		},
		c.Request().Context(),
	)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

    _, err = models.CreateUserToProduct(
        &models.UserToProduct{
            ProductId: product.Id,
            UserId: productJSON.UserId,
        },
        c.Request().Context(),
    )
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	_, err = models.CreateProductSetting(
		&models.ProductSetting{
			Active:    true,
			ProductId: product.Id,
			UserId:    productJSON.UserId,
		},
		c.Request().Context(),
	)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	priceFlt, err := strconv.ParseFloat(productJSON.Price[0].Amount, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product price")
		// TODO: Send DP
	}

	priceDTO := models.Price{
		Amount:    priceFlt,
		Currency:  productJSON.Price[0].Currency,
		ProductId: product.Id,
	}

	_, err = models.CreatePrice(priceDTO, c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product price")
		// TODO: Send DP
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product id")
		// Send DP
	}

	productId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// Send DP
	}

	productJSON := models.Product{
		Id: productId,
	}

	if err := c.Bind(&productJSON); err != nil {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product data")
		// TODO: Send DP
	}

	product, err := models.UpdateProduct(&productJSON, c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product id")
		// TODO: Send DP
	}

	productId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// Send DP
	}

	if err := models.DeleteProduct(productId, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	return c.NoContent(http.StatusOK)
}
