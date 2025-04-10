package group

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	_ "unsafe"
)

// 解决跨域问题
//
// 参考: https://blog.csdn.net/u011866450/article/details/126958238
func CORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent) // 禁止所有 OPTIONS 方法 原因见博文
	}
}

//go:linkname walk
func walk(path string) (files []string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		name := entry.Name()
		file := filepath.Join(path, name)
		if entry.IsDir() {
			files = append(files, walk(file)...)
		} else {
			files = append(files, file)
		}
	}
	return
}

// 绑定静态资源
func Static(s string) gin.HandlerFunc {
	s = filepath.Clean(s)
	repl := strings.NewReplacer(s, "", "\\", "/")
	files := make(map[string]string)
	for _, file := range walk(s) {
		files[repl.Replace(file)] = file
	}
	if index, ok := files["/index.html"]; ok {
		files["/"] = index
	}
	return func(ctx *gin.Context) {
		if file, ok := files[ctx.Request.URL.Path]; ok {
			ctx.File(file)
			ctx.Abort()
		}
	}
}
