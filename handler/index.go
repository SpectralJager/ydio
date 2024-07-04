package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/SpectralJager/ydio/service"
	"github.com/SpectralJager/ydio/view"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	Searcher *service.DownloadAudio
}

func (h IndexHandler) RenderPage(ctx echo.Context) error {
	sess, err := session.Get("session", ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}
	err = sess.Save(ctx.Request(), ctx.Response())
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return view.IndexView().Render(context.TODO(), ctx.Response())
}

func (h IndexHandler) SearchVideo(ctx echo.Context) error {
	url := ctx.QueryParam("url")
	meta, err := h.Searcher.GetAudioMetadate(url)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	err = SetValueToSession(ctx, "audioID", meta.ID)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, "/audio")
}

func (h IndexHandler) SearchPlaylist(ctx echo.Context) error {
	url := ctx.QueryParam("url")
	meta, err := h.Searcher.GetPlaylistMetadate(url)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	err = SetValueToSession(ctx, "playlistID", meta.ID)
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, "/playlist")
}
