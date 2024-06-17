package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SpectralJager/ydio/service"
	"github.com/SpectralJager/ydio/view"
	"github.com/labstack/echo/v4"
)

type AudioHandler struct {
	Downloader *service.DownloadAudio
}

func (h AudioHandler) RenderPage(ctx echo.Context) error {
	id := ctx.Param("id")
	meta, err := h.Downloader.GetAudioMetadate(id)
	if err != nil {
		ctx.Response().Header().Set("HX-Location", "/")
		return nil
	}
	return view.AudioView(meta).Render(context.TODO(), ctx.Response())
}

func (h AudioHandler) DownloadHTMX(ctx echo.Context) error {
	id := ctx.FormValue("id")
	meta, err := h.Downloader.GetAudioMetadate(id)
	if err != nil {
		return view.SearchResult(err.Error(), true).Render(context.TODO(), ctx.Response())
	}
	_, err = h.Downloader.DownloadAudio(meta)
	if err != nil {
		return view.SearchResult(err.Error(), true).Render(context.TODO(), ctx.Response())
	}
	ctx.Response().Header().Set("HX-Redirect", fmt.Sprintf("/audio/%s/get", id))
	return nil
}

func (h AudioHandler) GetAudio(ctx echo.Context) error {
	id := ctx.Param("id")
	meta, err := h.Downloader.GetAudioMetadate(id)
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	ctx.Response().Header().Set("Content-Type", "audio/mpeg, audio/x-mpeg, audio/mp3, audio/x-mp3, audio/mpeg3, audio/x-mpeg3, audio/mpg, audio/x-mpg, audio/x-mpegaudio")
	return ctx.Attachment(
		fmt.Sprintf("./public/audio/%s.mp3", id),
		fmt.Sprintf("%s.mp3", meta.Title),
	)
}
