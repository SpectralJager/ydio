package handler

import (
	"context"
	"fmt"

	"github.com/SpectralJager/ydio/service"
	"github.com/SpectralJager/ydio/view"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	Searcher *service.DownloadAudio
}

func (h IndexHandler) RenderPage(ctx echo.Context) error {
	return view.IndexView().Render(context.TODO(), ctx.Response())
}

func (h IndexHandler) SearchHTMX(ctx echo.Context) error {
	url := ctx.FormValue("url")
	meta, err := h.Searcher.GetAudioMetadate(url)
	if err != nil {
		return view.SearchResult(err.Error(), true).Render(context.TODO(), ctx.Response())
	}
	ctx.Response().Header().Set("HX-Location", fmt.Sprintf("/audio/%s", meta.ID))
	return nil
}
