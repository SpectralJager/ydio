package main

import (
	"log"

	"github.com/SpectralJager/ydio/view"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	// static files route
	app.Static("/static", "./public")

	// index view routes
	// `/`
	// TODO: implement index view routes
	app.GET("/", IndexHandle)

	// single audio download view routes
	// `/single/...`
	// TODO: implement index view routes
	app.GET("/single", SingleHandle)
	app.GET("/single/:id", SingleInfoHandle)
	app.GET("/single/:id/download", SingleDownloadHandle)

	// multiple audios download view routes
	// `/multiple/...`
	// TODO: implement multiple audios download view routes

	// playlist audio download view routes
	// `/playlist/...`
	// TODO: implement playlist audio download view routes

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func IndexHandle(c echo.Context) error {
	return view.IndexView().Render(c.Request().Context(), c.Response())
}

func SingleHandle(c echo.Context) error {
	return view.SingleView().Render(c.Request().Context(), c.Response())
}

func SingleInfoHandle(c echo.Context) error {
	return view.SingleInfoView().Render(c.Request().Context(), c.Response())
}

func SingleDownloadHandle(c echo.Context) error {
	return view.SingleDownloadView().Render(c.Request().Context(), c.Response())
}
