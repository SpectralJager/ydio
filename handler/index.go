package handler

import (
	"context"
	"fmt"
	"net/http"

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

func (h IndexHandler) SearchVideo(ctx echo.Context) error {
	url := ctx.QueryParam("url")
	meta, err := h.Searcher.GetAudioMetadate(url)
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/audio/%s", meta.ID))
}

func (h IndexHandler) SearchPlaylist(ctx echo.Context) error {
	url := ctx.QueryParam("url")
	meta, err := h.Searcher.GetPlaylistMetadate(url)
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/playlist/%s", meta.ID))
}
