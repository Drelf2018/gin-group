package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) (data any, err Error)

type Chain []HandlerFunc

// Call is a gin.HandlerFunc with receiver f.
func (f HandlerFunc) Handler(ctx *gin.Context) {
	if data, err := f(ctx); data != nil || err != nil {
		ctx.JSON(http.StatusOK, SetResponse(ctx, data, err))
	}
}
