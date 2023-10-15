package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Nicked/models"
)

type ProductJSON struct {
	Currency string
	Email    string
	Name     string
	Price    string
	Sku      string
	Store    string
	Url      string
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

func GetProductBySku(c echo.Context) error {
	sku := c.QueryParam("sku")
	if sku == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product sku")
		// TODO: Send DP
	}

	store := c.QueryParam("store")
	if store == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product store")
		// TODO: Send DP
	}

	email := c.QueryParam("email")
	if email == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid user email")
		// TODO: Send DP
	}

	user, err := models.GetUserByEmail(email, c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing email")
		// TODO: Send DP
	}

	if user != nil {
		product, err := models.GetProductBySku(sku, store, user.Id, c.Request().Context())
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
			// TODO: Send DP
		}

		if product != nil {
			return c.JSON(http.StatusOK, product)
		}
	}

	return c.NoContent(http.StatusNotFound)
}

func CreateProduct(c echo.Context) error {
	var productJSON ProductJSON
	if err := c.Bind(&productJSON); err != nil {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product")
		// TODO: Send DP
	}

	user, err := models.GetUserByEmail(productJSON.Email, c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
	}

	productDTO := models.Product{
		Name:   productJSON.Name,
		Active: true,
		Sku:    productJSON.Sku,
		Store:  productJSON.Store,
		Url:    productJSON.Url,
		UserId: user.Id,
	}

	product, err := models.CreateProduct(&productDTO, c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	priceFlt, err := strconv.ParseFloat(productJSON.Price, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product price")
		// TODO: Send DP
	}

	priceDTO := models.Price{
		Amount:    priceFlt,
		Currency:  productJSON.Currency,
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
	sku := c.Param("sku")
	if sku == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product id")
		// TODO: Send DP
	}

	store := c.QueryParam("store")
	if store == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product store")
		// TODO: Send DP
	}

	email := c.QueryParam("email")
	if email == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid user email")
		// TODO: Send DP
	}

	user, err := models.GetUserByEmail(email, c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing email")
		// TODO: Send DP
	}

	if user != nil {
		var productJSON models.Product
		if err := c.Bind(&productJSON); err != nil {
			return echo.NewHTTPError(http.StatusFailedDependency, "invalid product data")
			// TODO: Send DP
		}
		productJSON.Sku = sku
		productJSON.Store = store
		productJSON.UserId = user.Id

		product, err := models.UpdateProduct(&productJSON, c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
			// TODO: Send DP
		}

		return c.JSON(http.StatusOK, product)
	}

	return c.NoContent(http.StatusNotFound)
}

func DeleteProduct(c echo.Context) error {
	sku := c.Param("sku")
	if sku == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product id")
		// TODO: Send DP
	}

	store := c.QueryParam("store")
	if store == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product store")
		// TODO: Send DP
	}

	email := c.QueryParam("email")
	if email == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid user email")
		// TODO: Send DP
	}

	user, err := models.GetUserByEmail(email, c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing email")
		// TODO: Send DP
	}

	if user != nil {
		product, err := models.GetProductBySku(sku, store, user.Id, c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
			// TODO: Send DP
		}

		if err := models.DeleteProduct(product.Id, c.Request().Context()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
			// TODO: Send DP
		}

		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusNotFound)
}
