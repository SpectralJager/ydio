package main

import (
	"log"
	"os"

	"github.com/SpectralJager/ydio/handler"
	"github.com/SpectralJager/ydio/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()

	app.StaticFS("/static", os.DirFS("./public"))

	app.Use(middleware.Logger())

	audioService := service.NewDownloadAudioService()

	indexHandler := handler.IndexHandler{Searcher: audioService}
	audioHandler := handler.AudioHandler{Downloader: audioService}

	htmx := app.Group("/htmx/v1")
	htmx.POST("/search", indexHandler.SearchHTMX)
	htmx.POST("/download", audioHandler.DownloadHTMX)

	index := app.Group("/")
	index.GET("", indexHandler.RenderPage)

	audio := app.Group("/audio")
	audio.GET("/:id", audioHandler.RenderPage)
	audio.GET("/:id/get", audioHandler.GetAudio)

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
