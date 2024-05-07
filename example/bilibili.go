package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	group "github.com/Drelf2018/gin-group"
	"github.com/gin-gonic/gin"
)

type View struct {
	Data struct {
		Pic string
	}
}

var ErrInvalidBvid = group.NewError(1, "example: %v is an invalid bvid")

func GetPicColonmid(ctx *gin.Context) (any, group.Error) {
	mid := ctx.Param("mid")

	if !strings.HasPrefix(mid, "BV") {
		return nil, ErrInvalidBvid.Format(mid)
	}

	resp, err := http.Get(fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", mid))
	if err != nil {
		return nil, group.AutoError(err)
	}

	p, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, group.AutoError(err)
	}

	var view View
	err = json.Unmarshal(p, &view)
	if err != nil {
		return nil, group.AutoError(err)
	}

	return view.Data.Pic, nil
}
