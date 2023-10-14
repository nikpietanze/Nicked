package handlers

import (
	"Nicked/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateDataPoint(c echo.Context) error {
	var dataPointJSON models.DataPoint
	if err := c.Bind(&dataPointJSON); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "missing or invalid data point")
	}

    if err := models.CreateDataPoint(&dataPointJSON, c.Request().Context()); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

    return c.NoContent(200)
}
