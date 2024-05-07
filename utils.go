package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 解决跨域问题
//
// 参考: https://blog.csdn.net/u011866450/article/details/126958238
func CORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	// 禁止所有 OPTIONS 方法 原因见博文
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
}

const UserKey string = "__user__"

// GetUser returns the user as type U.
func GetUser[U any](ctx *gin.Context) (u U) {
	if val, ok := ctx.Get(UserKey); ok && val != nil {
		u, _ = val.(U)
	}
	return
}

func SetUser(ctx *gin.Context, user any) {
	ctx.Set(UserKey, user)
}

func Abort(ctx *gin.Context, data any, err Error) {
	ctx.AbortWithStatusJSON(http.StatusOK, SetResponse(ctx, data, err))
}
