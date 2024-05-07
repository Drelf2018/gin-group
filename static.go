package group

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func Static(s string) gin.HandlerFunc {
	s = filepath.Clean(s)
	files := make(map[string]string)
	rep := strings.NewReplacer(s, "", "\\", "/")

	for _, file := range walk(s) {
		files[rep.Replace(file)] = file
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
