package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"Nicked/db"
	"Nicked/handlers"
	apiHandlers "Nicked/handlers/api"
	"Nicked/middlewares"
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

	templates := make(map[string]*template.Template)
	templates["home.html"] = template.Must(template.ParseFiles("views/home.html", "views/layouts/base.html"))
	templates["privacy.html"] = template.Must(template.ParseFiles("views/privacy.html", "views/layouts/base.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.Static("/static", "public")

	// global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	scraperStarted := false
	e.GET("/", func(c echo.Context) error {
		if !scraperStarted {
			scraper.Init()
			scraperStarted = true
		}

		return c.Render(http.StatusOK, "home.html", map[string]interface{}{
			"title": "Nicked",
		})
	})

	e.GET("/privacy", func(c echo.Context) error {
		return c.Render(http.StatusOK, "privacy.html", map[string]interface{}{
			"title": "Privacy Policy | Nicked",
		})
	})

	api := e.Group("/api")
	{
		// api middleware
		api.Use(middlewares.Auth())

		api.POST("/analytics", handlers.CreateDataPoint)

		user := api.Group("/user")
		{
			user.GET("/{id}", apiHandlers.GetUser)
			user.POST("/", apiHandlers.CreateUser)
			user.PUT("/", apiHandlers.UpdateUser)
			user.DELETE("/{id}", apiHandlers.DeleteUser)
		}

		item := api.Group("/item")
		{
			item.GET("/{id}", apiHandlers.GetItem)
			item.POST("/", apiHandlers.CreateItem)
			item.PUT("/", apiHandlers.UpdateItem)
			item.DELETE("/{id}", apiHandlers.DeleteItem)
		}

		price := api.Group("/price")
		{
			price.GET("/{id}", apiHandlers.GetPrice)
			price.POST("/", apiHandlers.CreatePrice)
			price.DELETE("/{id}", apiHandlers.DeletePrice)
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
