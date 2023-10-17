package main

import (
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"Nicked/db"
	"Nicked/handlers"
	apiHandlers "Nicked/handlers/api"
	"Nicked/middlewares"
	"Nicked/models"
	"Nicked/scraper"
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
	initTables(context.Background())

    // Templates
	templates := make(map[string]*template.Template)
	templates["home.html"] = template.Must(template.ParseFiles("views/home.html", "views/layouts/base.html"))
	templates["privacy.html"] = template.Must(template.ParseFiles("views/privacy.html", "views/layouts/base.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

    // Static Files
	e.Static("/static", "public")

	// Global Middleware
	e.Use(middleware.Logger())
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

    // /api/prices
    api.GET("/price/:id", apiHandlers.GetPrice)
	api.POST("/price", apiHandlers.CreatePrice)
	api.DELETE("/price/:id", apiHandlers.DeletePrice)

    // Scraper
	scraper.Init()

	e.Logger.Fatal(e.Start(":8080"))
}

func initTables(ctx context.Context) {
	if err := models.InitAnalytics(ctx); err != nil {
		log.Println(err)
	}
	if err := models.InitUser(ctx); err != nil {
		log.Println(err)
	}
	if err := models.InitProduct(ctx); err != nil {
		log.Println(err)
	}
	if err := models.InitPrice(ctx); err != nil {
		log.Println(err)
	}
}
