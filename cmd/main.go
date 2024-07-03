package main

import (
	"log"
	"os"

	"github.com/SpectralJager/ydio/handler"
	"github.com/SpectralJager/ydio/service"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	index.GET("/searchVideo", indexHandler.SearchVideo)
	index.GET("/searchPlaylist", indexHandler.SearchPlaylist)

	audio := app.Group("/audio/:id")
	audio.GET("", audioHandler.RenderPage)
	audio.GET("/download", audioHandler.DownloadAudio)
	audio.GET("/get", audioHandler.GetAudio)
	audio.GET("/status", audioHandler.GetStatus)

	playlist := app.Group("/playlist/:id", session.Middleware(sessions.NewCookieStore([]byte("test"))))
	playlist.GET("", playlistHandler.RenderPage)
	playlist.POST("/download", playlistHandler.DownloadPlaylist)
	playlist.GET("/get", playlistHandler.GetPlaylist)
	playlist.GET("/status", playlistHandler.StatusPlaylist)

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
