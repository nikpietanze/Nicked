package middlewares

import (
	"Nicked/config"
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Auth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(config.AUTH_USERNAME)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(config.AUTH_PASSWORD)) == 1 {
			return true, nil
		}
		return false, nil
	})
}
