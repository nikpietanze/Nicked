package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Nicked/models"
)

func GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user id")
        // Send DP
	}

	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	user, err := models.GetUser(&userId, c.Request().Context())
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	return c.JSON(http.StatusOK, user)
}

func GetUserByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user email")
        // Send DP
	}

	user, err := models.GetUserByEmail(email, c.Request().Context())
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var userJSON models.User
	if err := c.Bind(&userJSON); err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user data")
        // Send DP
	}

	user, err := models.CreateUser(&userJSON, c.Request().Context())
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
    id := c.Param("id")
	if id == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user id")
        // Send DP
	}

	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    userJSON := models.User{
        Id: userId,
    }

	if err := c.Bind(&userJSON); err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user data")
        // Send DP
	}

	user, err := models.UpdateUser(&userJSON, c.Request().Context())
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
        return echo.NewHTTPError(http.StatusFailedDependency, "invalid user id")
        // Send DP
	}

	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

	if err := models.DeleteUser(userId, c.Request().Context()); err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError, "error processing user")
        // Send DP
	}

    return c.NoContent(http.StatusOK)
}
