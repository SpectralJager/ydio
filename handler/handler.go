package handler

import (
	"log"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetValueFromSession[T any](ctx echo.Context, key string) (T, bool) {
	var value T
	var ok bool
	sess, err := session.Get("session", ctx)
	if err != nil {
		log.Println(err)
		return value, false
	}
	value, ok = sess.Values[key].(T)
	return value, ok
}

func SetValueToSession[T any](ctx echo.Context, key string, value T) error {
	sess, err := session.Get("session", ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	sess.Values[key] = value
	err = sess.Save(ctx.Request(), ctx.Response())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
