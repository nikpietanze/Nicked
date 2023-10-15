package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Nicked/models"
)

func GetPrice(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid price id")
        // Send DP
	}

	priceId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing price")
        // Send DP
	}

	price, err := models.GetPrice(priceId, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing price")
        // Send DP
	}

	return c.JSON(http.StatusOK, price)
}

func CreatePrice(c echo.Context) error {
	var priceJSON models.Price
	if err := c.Bind(&priceJSON); err != nil {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid price")
        // Send DP
	}

	price, err := models.CreatePrice(priceJSON, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing price")
        // Send DP
	}

    return c.JSON(http.StatusOK, price)
}

func DeletePrice(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid price id")
        // Send DP
	}

	priceId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing price")
        // Send DP
	}

	if err := models.DeletePrice(priceId, c.Request().Context()); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing price")
        // Send DP
	}

    return c.NoContent(http.StatusOK)
}
