package middlewares

import (
	"Nicked/config"
	"strings"

	"github.com/kataras/iris/v12"
)

func Auth() iris.Handler {
	return func(ctx iris.Context) {
		pathParts := strings.Split(ctx.Path(), "/")
		if pathParts[0] != "api" {
			ctx.Next()
		}

		authenticated := false
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			if authHeader == config.NICKED_EXT_API_KEY {
				authenticated = true
			}
		}
		if !authenticated {
			ctx.StopWithText(iris.StatusUnauthorized, "not authorized")
		}
		ctx.Next()
	}
}
