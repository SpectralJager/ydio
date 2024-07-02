package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SpectralJager/ydio/service"
	"github.com/SpectralJager/ydio/view"
	"github.com/labstack/echo/v4"
)

type PlaylistHandler struct {
	Downloader *service.DownloadAudio
}

func (h PlaylistHandler) RenderPage(ctx echo.Context) error {
	id := ctx.Param("id")
	meta, err := h.Downloader.GetPlaylistMetadate(id)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	var totalDuration time.Duration
	for _, video := range meta.Videos {
		totalDuration += video.Duration
	}
	return view.PlaylistView(meta, totalDuration).Render(context.TODO(), ctx.Response())
}

func (h PlaylistHandler) DownloadPlaylist(ctx echo.Context) error {
	id := ctx.Param("id")
	meta, err := h.Downloader.GetPlaylistMetadate(id)
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	ids, ok := ctx.QueryParams()["videoIds[]"]
	if !ok {
		ids = []string{}
	}
	return view.PlaylistDownload(meta, ids).Render(context.TODO(), ctx.Response())
}

func (h PlaylistHandler) GetPlaylist(ctx echo.Context) error {
	id := ctx.Param("id")
	meta, err := h.Downloader.GetPlaylistMetadate(id)
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	ctx.Response().Header().Set("Content-Type", "application/zip")
	return ctx.Attachment(
		fmt.Sprintf("./public/audio/%s.zip", id),
		fmt.Sprintf("%s.zip", meta.Title),
	)
}

func (h PlaylistHandler) StatusPlaylist(ctx echo.Context) error {
	id := ctx.Param("id")
	ids, ok := ctx.QueryParams()["videoIds[]"]
	if !ok {
		ids = []string{}
	}
	log.Println(ids)
	meta, err := h.Downloader.GetPlaylistMetadate(id)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	err = h.Downloader.DownloadPlaylist(meta, ids)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	w := ctx.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Println("sse")

	var buff bytes.Buffer
	buff.WriteString("event: close\n")
	buff.WriteString("data: ")
	view.PlaylistGet(meta).Render(context.TODO(), &buff)
	buff.WriteString("\n\n")

	w.Write(buff.Bytes())

	w.Flush()

	return view.Render(context.TODO(), ctx.Response(), view.PlaylistGet(meta))
}
