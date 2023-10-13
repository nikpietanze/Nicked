package middlewares

import (
	"log"
	"time"

	"github.com/kataras/iris/v12"
)

func Logger() iris.Handler {
    return func(ctx iris.Context) {
        t := time.Now()

        ctx.Next()

        path:= ctx.Path()
        latency := time.Since(t)
        status := ctx.GetStatusCode()
        log.Println(path, status, latency)
    }
}
