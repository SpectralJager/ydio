package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SpectralJager/ydio/service"
	"github.com/SpectralJager/ydio/view"
	"github.com/labstack/echo/v4"
)

type AudioHandler struct {
	Downloader *service.DownloadAudio
}

func (h AudioHandler) RenderPage(ctx echo.Context) error {
	id, ok := GetValueFromSession[string](ctx, "audioID")
	if !ok {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	meta, err := h.Downloader.GetAudioMetadate(id)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return view.AudioView(meta).Render(context.TODO(), ctx.Response())
}

func (h AudioHandler) DownloadAudio(ctx echo.Context) error {
	return view.Render(context.TODO(), ctx.Response(), view.AudioStartDownload())
}

func (h AudioHandler) GetAudio(ctx echo.Context) error {
	id, ok := GetValueFromSession[string](ctx, "audioID")
	if !ok {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	meta, err := h.Downloader.GetAudioMetadate(id)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	ctx.Response().Header().Set("Content-Type", "audio/mpeg, audio/x-mpeg, audio/mp3, audio/x-mp3, audio/mpeg3, audio/x-mpeg3, audio/mpg, audio/x-mpg, audio/x-mpegaudio")
	return ctx.Attachment(
		fmt.Sprintf("./public/audio/%s.mp3", id),
		fmt.Sprintf("%s.mp3", meta.Title),
	)
}

func (h AudioHandler) GetStatus(ctx echo.Context) error {
	w := ctx.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id, ok := GetValueFromSession[string](ctx, "audioID")
	if !ok {
		SendClose(ctx)
		return nil
	}
	meta, err := h.Downloader.GetAudioMetadate(id)
	if err != nil {
		log.Println(err)
		SendClose(ctx)
		return nil
	}
	err = h.Downloader.DownloadAudio(meta)
	if err != nil {
		log.Println(err, 1)
		SendClose(ctx)
		return nil
	}

	var buff bytes.Buffer
	buff.WriteString("event: close\n")
	buff.WriteString("data: ")
	view.AudioGet().Render(context.TODO(), &buff)
	buff.WriteString("\n\n")

	w.Write(buff.Bytes())

	w.Flush()

	return nil
}
