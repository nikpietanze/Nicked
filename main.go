package main

import (
	"github.com/kataras/iris/v12"

	"Nicked/db"
	"Nicked/handlers"
	apiHandlers "Nicked/handlers/api"
	"Nicked/middlewares"
	"Nicked/models"
	"Nicked/scraper"
)

func main() {
	app := iris.New()

	db.Init()

	app.RegisterView(iris.HTML("./views", ".html"))
	app.HandleDir("/public", iris.Dir("./public"))

    // middlewares
	app.Use(iris.Compression)
    app.Use(middlewares.Auth())

	scraperStarted := false
	app.Get("/", func(ctx iris.Context) {
		if !scraperStarted {
			scraper.Init(ctx)
			scraperStarted = true
		}

		data := iris.Map{
			"Title": "Home | Nicked",
		}

		ctx.ViewLayout("layouts/main")

		err := ctx.View("index", data)
		if err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
	})

	app.Get("/privacy", func(ctx iris.Context) {
		data := iris.Map{
			"Title": "Privacy | Nicked",
		}

		ctx.ViewLayout("layouts/main")

		err := ctx.View("privacy", data)
		if err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
	})

	app.Post("/analytics", handlers.CreateDataPoint)

	api := app.Party("/api")
	{
		user := api.Party("/user")
		{
			user.Get("/{id}", apiHandlers.GetUser)
			user.Post("/", apiHandlers.CreateUser)
			user.Put("/", apiHandlers.UpdateUser)
			user.Delete("/{id}", apiHandlers.DeleteUser)
		}

		item := api.Party("/item")
		{
			item.Get("/{id}", apiHandlers.GetItem)
			item.Post("/", apiHandlers.CreateItem)
			item.Put("/", apiHandlers.UpdateItem)
			item.Delete("/{id}", apiHandlers.DeleteItem)
		}

		price := api.Party("/price")
		{
			price.Get("/{id}", apiHandlers.GetPrice)
			price.Post("/", apiHandlers.CreatePrice)
			price.Delete("/{id}", apiHandlers.DeletePrice)
		}
	}

	app.Listen(":8080")
}
