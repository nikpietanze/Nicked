package main

import (
	"context"
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nicked.io/db"
	"nicked.io/handlers"
	apiHandlers "nicked.io/handlers/api"
	"nicked.io/middlewares"
	"nicked.io/models"
	"nicked.io/scraper"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found: " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

func main() {
	e := echo.New()

	db.Init()
	db.Client.RegisterModel((*models.UserToProduct)(nil))

	err := createSchema(context.Background())
	if err != nil {
		panic(err)
	}

	// Templates
	templates := make(map[string]*template.Template)
	templates["home.html"] = template.Must(template.ParseFiles("views/home.html", "views/layouts/base.html"))
	templates["privacy.html"] = template.Must(template.ParseFiles("views/privacy.html", "views/layouts/base.html"))
	templates["info/sales-and-discounts-tracking/v1.html"] = template.Must(template.ParseFiles("views/info/sales-and-discounts-tracking/v1.html", "views/layouts/base.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	// Static Files
	e.Static("/static", "public")

	// Global Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// Website Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", map[string]interface{}{
			"title": "Nicked",
		})
	})
	e.GET("/privacy", func(c echo.Context) error {
		return c.Render(http.StatusOK, "privacy.html", map[string]interface{}{
			"title": "Privacy Policy | Nicked",
		})
	})

	// Info Routes
	info := e.Group("/info")
	info.GET("/sales-and-discounts-tracking", func(c echo.Context) error {
		return c.Render(http.StatusOK, "info/sales-and-discounts-tracking/v1.html", map[string]interface{}{
			"title":       "Save More, Shop Smarter: The Power of a Sales and Discount Tracker!",
			"description": "Discover the Ultimate Sales and Discount Tracker for Amazon Shoppers! Maximize Savings, Stay Organized, and Never Miss a Deal. Your Must-Have Companion for a Smarter Amazon Shopping Experience. Start Supercharging Your Savings Today!",
		})
	})

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
	scraper.Init()

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
