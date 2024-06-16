package handler

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/SpectralJager/ydio/view"
	"github.com/kkdai/youtube/v2"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) RenderPage(ctx echo.Context) error {
	return view.IndexView().Render(context.TODO(), ctx.Response())
}

func (h IndexHandler) DownloadFile(ctx echo.Context) error {
	return ctx.Attachment("audio.webm", "audio.webm")
}

func (h IndexHandler) SearchVideo(ctx echo.Context) error {
	videoURL := ctx.FormValue("video")
	client := youtube.Client{
		MaxRoutines: 2,
	}
	meta, err := client.GetVideo(videoURL)
	if err != nil {
		return view.SearchVideoResponse(echo.Map{
			"error": err.Error(),
			"url":   videoURL,
		}).Render(context.TODO(), ctx.Response())
	}
	formats := meta.Formats.WithAudioChannels()
	formats = formats.Select(func(f youtube.Format) bool {
		return strings.Contains(f.MimeType, "audio/webm")
	})
	stream, _, err := client.GetStream(meta, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create("audio.webm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	ctx.Response().Header().Set("HX-Redirect", "/audio")
	return nil
}
