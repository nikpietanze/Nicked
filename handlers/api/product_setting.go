package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Nicked/models"
)

func UpdateProductSetting(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product id")
		// Send DP
	}

	productSettingId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// Send DP
	}

	productSettingJSON := models.ProductSetting{
        Id: productSettingId,
    };

	if err := c.Bind(&productSettingJSON); err != nil {
		return echo.NewHTTPError(http.StatusFailedDependency, "invalid product data")
		// TODO: Send DP
	}

	productSetting, err := models.UpdateProductSetting(&productSettingJSON, c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error processing product")
		// TODO: Send DP
	}

	return c.JSON(http.StatusOK, productSetting)
}
