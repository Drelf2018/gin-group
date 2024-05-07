package main

import (
	"time"

	group "github.com/Drelf2018/gin-group"
	"github.com/gin-gonic/gin"
)

type Response struct {
	group.DefaultResponse
	Time string `json:"time"`
	IP   string `json:"ip"`
}

func (r *Response) Set(ctx *gin.Context, data any, err group.Error) {
	r.IP = ctx.ClientIP()
	r.DefaultResponse.Set(ctx, data, err)
}

func init() {
	group.NewResponse = func() group.Response {
		return &Response{
			Time: time.Now().Format(time.DateTime),
		}
	}
}
