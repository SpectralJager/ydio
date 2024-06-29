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
	playlistHandler := handler.PlaylistHandler{Downloader: audioService}

	index := app.Group("")
	index.GET("", indexHandler.RenderPage)
	index.GET("/search", indexHandler.SearchHTMX)

	audio := app.Group("/audio/:id")
	audio.GET("", audioHandler.RenderPage)
	audio.GET("/download", audioHandler.DownloadAudio)
	audio.GET("/get", audioHandler.GetAudio)
	audio.GET("/event", audioHandler.GetStatus)

	playlist := app.Group("/playlist/:id")
	playlist.GET("", playlistHandler.RenderPage)
	playlist.GET("/download", playlistHandler.DownloadPlaylist)
	playlist.GET("/get", playlistHandler.GetPlaylist)

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
