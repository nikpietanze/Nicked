package main

import (
	"github.com/kataras/iris/v12"

	"Nicked/db"
	"Nicked/handlers"
	"Nicked/models"
	"Nicked/scraper"
)

func main() {
	app := iris.New()
	app.Use(iris.Compression)

	db.Init()

    app.Get("/", middleware)

	api := app.Party("/api")
	{
		user := api.Party("/user")
		{
			user.Get("/{id}", func(ctx iris.Context) {
				println("get user endpoint hit")
				models.InitUser(ctx)
				handlers.GetUser(ctx)
			})
			user.Post("/", handlers.CreateUser)
			user.Put("/", handlers.UpdateUser)
			user.Delete("/{id}", handlers.DeleteUser)
		}

		item := api.Party("/item")
		{
			item.Get("/{id}", func(ctx iris.Context) {
                println("get item endpoint hit")
				models.InitItem(ctx)
				handlers.GetItem(ctx)
			})
			item.Post("/", handlers.CreateItem)
			item.Put("/", handlers.UpdateItem)
			item.Delete("/{id}", handlers.DeleteItem)
		}

		price := api.Party("/price")
		{
			price.Get("/{id}", func(ctx iris.Context) {
				println("get price endpoint hit")
				models.InitPrice(ctx)
				handlers.GetPrice(ctx)
			})
			price.Post("/", handlers.CreatePrice)
			price.Delete("/{id}", handlers.DeletePrice)
		}
	}


	app.Listen(":8080")
}

func middleware(ctx iris.Context) {
    scraper.Init(ctx)
}
