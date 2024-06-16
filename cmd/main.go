package main

import (
	"log"
	"os"

	"github.com/SpectralJager/ydio/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	app.StaticFS("/static", os.DirFS("./public"))

	indexHandler := handler.IndexHandler{}

	app.GET("/", indexHandler.RenderPage)
	app.GET("/audio", indexHandler.DownloadFile)
	app.POST("/htmx/search_video", indexHandler.SearchVideo)

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}

}
