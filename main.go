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
        tracker := api.Party("/tracker")
        {
            tracker.Get("/{id}", handlers.GetTracker)
            tracker.Post("/", handlers.CreateTracker)
            tracker.Put("/", handlers.UpdateTracker)
            tracker.Delete("/{id}", handlers.DeleteTracker)
        }
    }

    scraper.Init()

    app.Listen(":8080")
}

