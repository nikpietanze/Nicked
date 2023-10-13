package middlewares

import (
	"github.com/kataras/iris/v12"
)

func Cors() iris.Handler {
    return func(ctx iris.Context) {
        ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
            ctx.Header("Access-Control-Methods",
                "POST, PUT, PATCH, DELETE")

            ctx.Header("Access-Control-Allow-Headers",
                "Access-Control-Allow-Origin,Content-Type,Authorization")

            ctx.Header("Access-Control-Max-Age",
                "86400")

			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
    }
}
