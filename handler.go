package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) (data any, err error)

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

// 默认转换器
var DefaultConvertor = func(f HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if data, err := f(ctx); err == nil {
			if data != nil {
				ctx.JSON(http.StatusOK, Response{0, "", data})
			}
		} else {
			if code, ok := data.(int); ok {
				ctx.JSON(http.StatusOK, Response{code, err.Error(), nil})
			} else {
				ctx.JSON(http.StatusOK, Response{1, err.Error(), data})
			}
		}
	}
}
