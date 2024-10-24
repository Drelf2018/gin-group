package group

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Group struct {
	Path        string
	Middleware  gin.HandlerFunc
	Middlewares []gin.HandlerFunc
	Custom      func(gin.IRouter)
	Handlers    []HandlerFunc
	Groups      []Group
}

func (group *Group) Bind(r gin.IRouter) {
	if group.Path != "" {
		r = r.Group(group.Path)
	}
	if group.Middleware != nil {
		r.Use(group.Middleware)
	}
	if group.Middlewares != nil {
		r.Use(group.Middlewares...)
	}
	if group.Custom != nil {
		group.Custom(r)
	}
	for _, handler := range group.Handlers {
		var method, path string
		name := NameOfFunction(handler)
		if val, ok := pathCache.Load(name); ok {
			method = val.([2]string)[0]
			path = val.([2]string)[1]
		}
		if method == "" {
			matched := MethodExpr.FindStringSubmatch(name)
			if len(matched) == 3 {
				method = matched[1]
				path = ParsePath(matched[2])
			}
		}
		if method != "" {
			r.Handle(strings.ToUpper(method), path, handler.Handle)
		}
	}
	for _, v := range group.Groups {
		v.Bind(r)
	}
}

func (group *Group) New() (r *gin.Engine) {
	r = gin.New()
	group.Bind(r)
	return
}

func (group *Group) Default() (r *gin.Engine) {
	r = gin.Default()
	group.Bind(r)
	return
}
