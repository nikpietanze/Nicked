package main

import (
	"github.com/kataras/iris/v12"

	"pricetracker/db"
	"pricetracker/handlers"
	"pricetracker/scraper"
)

func main() {
    app := iris.New()
    app.Use(iris.Compression)

    db.Init()

    api := app.Party("/api")
    {
        user := api.Party("/user")
        {
            user.Get("/{id}", handlers.GetUser)
            user.Post("/", handlers.CreateUser)
            user.Put("/", handlers.UpdateUser)
            user.Delete("/{id}", handlers.DeleteUser)
        }

        item := api.Party("/item")
        {
            item.Get("/{id}", handlers.GetItem)
            item.Post("/", handlers.CreateItem)
            item.Put("/", handlers.UpdateItem)
            item.Delete("/{id}", handlers.DeleteItem)
        }
    }

    scraper.Init()

    app.Listen(":8080")
}

