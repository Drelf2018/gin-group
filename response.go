package group

import "github.com/gin-gonic/gin"

type Response interface {
	Set(ctx *gin.Context, data any, err Error)
}

type DefaultResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func (r *DefaultResponse) Set(ctx *gin.Context, data any, err Error) {
	r.Data = data
	if err != nil {
		r.Code = err.Code()
		r.Error = err.Error()
	}
}

var NewResponse = func() Response {
	return new(DefaultResponse)
}

func SetResponse(ctx *gin.Context, data any, err Error) (resp Response) {
	resp = NewResponse()
	resp.Set(ctx, data, err)
	return
}
