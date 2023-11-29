package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nicked.io/db"
	"nicked.io/handlers"
	apiHandlers "nicked.io/handlers/api"
	"nicked.io/middlewares"
	"nicked.io/models"
	"nicked.io/scraper"
)

func main() {
	e := echo.New()

	db.Init()
	db.Client.RegisterModel((*models.UserToProduct)(nil))

	err := createSchema(context.Background())
	if err != nil {
		panic(err)
	}

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// API Routes
	api := e.Group("/api")
	api.Use(middlewares.Auth())

	api.POST("/analytics", handlers.CreateDataPoint)

	// /api/users
	api.GET("/user/:id", apiHandlers.GetUser)
	api.GET("/user", apiHandlers.GetUserByEmail)
	api.POST("/user", apiHandlers.CreateUser)
	api.PUT("/user/:id", apiHandlers.UpdateUser)
	api.DELETE("/user/:id", apiHandlers.DeleteUser)

	// /api/products
	api.GET("/product/:id", apiHandlers.GetProduct)
	api.POST("/product", apiHandlers.CreateProduct)
	api.PUT("/product/:id", apiHandlers.UpdateProduct)
	api.DELETE("/product/:id", apiHandlers.DeleteProduct)

	// /api/products
	api.PUT("/product/:id", apiHandlers.UpdateProductSetting)

	// /api/prices
	api.GET("/price/:id", apiHandlers.GetPrice)
	api.POST("/price", apiHandlers.CreatePrice)
	api.DELETE("/price/:id", apiHandlers.DeletePrice)

	// Scraper
	go scraper.Init()

	e.Logger.Fatal(e.Start(":8080"))
}

func createSchema(ctx context.Context) error {
	models := []interface{}{
		(*models.User)(nil),
		(*models.Product)(nil),
		(*models.UserToProduct)(nil),
		(*models.ProductSetting)(nil),
		(*models.Price)(nil),
		(*models.DataPoint)(nil),
	}

	for _, model := range models {
		if _, err := db.Client.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
