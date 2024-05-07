package group

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Group struct {
	Path        string
	Middlewares gin.HandlersChain
	Customize   func(gin.IRouter)
	Handlers    Chain
	Groups      []Group
}

func (group *Group) Bind(r gin.IRouter) {
	if group.Path != "" {
		r = r.Group(group.Path)
	}

	r.Use(group.Middlewares...)

	if group.Customize != nil {
		group.Customize(r)
	}

	for _, handler := range group.Handlers {
		method, path := SplitHandlerName(handler)
		if method != "" {
			r.Handle(strings.ToUpper(method), ParsePath(path), handler.Handler)
		}
	}

	for _, v := range group.Groups {
		v.Bind(r)
	}
}

func New(group Group) (r *gin.Engine) {
	r = gin.New()
	group.Bind(r)
	return
}

func Default(group Group) (r *gin.Engine) {
	r = gin.Default()
	group.Bind(r)
	return
}
