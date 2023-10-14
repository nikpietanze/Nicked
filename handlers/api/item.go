package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Nicked/models"
)

type ItemJSON struct {
	Currency string
	Email    string
	Name     string
	Price    string
	Sku      string
	Store    string
	Url      string
}

func GetItem(c echo.Context) error {
	strId := c.QueryParam("id")
	if strId == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid item id")
        // TODO: Send DP
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

	item, err := models.GetItem(&id, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

	return c.JSON(http.StatusOK, item)
}

func CreateItem(c echo.Context) error {

	var itemJSON ItemJSON
	if err := c.Bind(&itemJSON); err != nil {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid item")
        // TODO: Send DP
	}

	user, err := models.GetUserByEmail(itemJSON.Email, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
	}

	itemDTO := models.Item{
		Name:     itemJSON.Name,
		IsActive: true,
		Sku:      itemJSON.Sku,
		Url:      itemJSON.Url,
		UserId:   user.Id,
	}

	item, err := models.CreateItem(&itemDTO, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

	priceFlt, err := strconv.ParseFloat(itemJSON.Price, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

	priceDTO := models.Price{
		Amount:   priceFlt,
		Currency: itemJSON.Currency,
		Store:    itemJSON.Store,
		ItemId:   item.Id,
	}

	_, err = models.CreatePrice(&priceDTO, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing price")
        // TODO: Send DP
	}

	return c.JSON(http.StatusOK, item)
}

func UpdateItem(c echo.Context) error {
	var itemJSON models.Item
	if err := c.Bind(&itemJSON); err != nil {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid item")
        // TODO: Send DP
	}

	item, err := models.UpdateItem(&itemJSON, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

	return c.JSON(http.StatusOK, item)
}

func DeleteItem(c echo.Context) error {
    strId := c.QueryParam("id")
	if strId == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid item")
        // TODO: Send DP
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

	if err := models.DeleteItem(id, c.Request().Context()); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing item")
        // TODO: Send DP
	}

    return c.NoContent(http.StatusOK)
}
