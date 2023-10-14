package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Nicked/models"
)

func GetUser(c echo.Context) error {
	strId := c.QueryParam("id")
	if strId == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user")
        // Send DP
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	user, getErr := models.GetUser(&id, c.Request().Context())
	if getErr != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var userJSON models.User
	if err := c.Bind(&userJSON); err != nil {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user")
        // Send DP
	}

	user, err := models.CreateUser(&userJSON, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	var userJSON models.User
	if err := c.Bind(&userJSON); err != nil {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user")
        // Send DP
	}

	user, err := models.UpdateUser(&userJSON, c.Request().Context())
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	strId := c.QueryParam("id")
	if strId == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user")
        // Send DP
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	if err := models.DeleteUser(id, c.Request().Context()); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    return c.NoContent(http.StatusOK)
}
