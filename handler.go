package group

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) (data any, err error)

type H = HandlerFunc

type Error struct {
	Code int
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func E(code int, err error) Error {
	return Error{code, err}
}

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

// Call is a gin.HandlerFunc with receiver f.
func (f HandlerFunc) Handle(ctx *gin.Context) {
	if data, err := f(ctx); err == nil {
		if data != nil {
			ctx.JSON(http.StatusOK, Response{0, "", data})
		}
	} else {
		ctx.Error(err)
		if e, ok := err.(Error); ok {
			ctx.JSON(http.StatusOK, Response{e.Code, e.Err.Error(), data})
		} else if code, ok := data.(int); ok {
			ctx.JSON(http.StatusOK, Response{code, err.Error(), nil})
		} else {
			ctx.JSON(http.StatusOK, Response{1, err.Error(), data})
		}
	}
}

type R[T any] struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  T      `json:"data"`
}

func (r R[T]) Unwrap() error {
	if r.Code == 0 {
		return nil
	}
	return fmt.Errorf(`group: ResponseError: "%s" with code %d`, r.Error, r.Code)
}

func Unmarshal[T any](data []byte) (r R[T], err error) {
	err = json.Unmarshal(data, &r)
	if err == nil {
		err = r.Unwrap()
	}
	return
}
