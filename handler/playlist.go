package handler

import (
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
	return ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/playlist/%s/get", meta.ID))
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
