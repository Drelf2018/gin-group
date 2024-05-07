package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	group "github.com/Drelf2018/gin-group"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct {
	MIME string
	Data []byte
}

var cache = make(map[string]string)
var files = make(map[string]File)

var ErrNoFile = errors.New("example: the file does not exist")

func GetResourceColonfile(ctx *gin.Context) (any, group.Error) {
	file, ok := files[ctx.Param("file")]
	if !ok {
		return nil, group.AutoError(ErrNoFile)
	}

	ctx.Data(http.StatusOK, file.MIME, file.Data)
	return nil, nil
}

func GetDownload(ctx *gin.Context) (any, group.Error) {
	url := ctx.Query("url")

	if file, ok := cache[url]; ok {
		return file, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, group.AutoError(err)
	}

	p, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, group.AutoError(err)
	}

	s := strings.ReplaceAll(uuid.New().String(), "-", "")
	cache[url] = s
	files[s] = File{
		MIME: resp.Header.Get("Content-Type"),
		Data: p,
	}
	return s, nil
}
