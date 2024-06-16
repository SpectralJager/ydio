package view

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(ctx context.Context, w http.ResponseWriter, component templ.Component) error {
	return component.Render(ctx, w)
}
