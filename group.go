package group

import (
	"github.com/gin-gonic/gin"
)

type (
	M = []gin.HandlerFunc
	H = []HandlerFunc
	G = []Group
)

// 接口组 (r = gin.IRouter)
type Group struct {
	// 相对路径 (r.Group)
	Path string

	// 中间件 (r.Use)
	Middlewares []gin.HandlerFunc

	// 自定义函数
	// 用户可以自行绑定内容
	CustomFunc func(r gin.IRouter)

	// 自动接口绑定 (r.Handle)
	Handlers []HandlerFunc

	// 自定义路径绑定 (r.Handle)
	HandlerMap map[string]HandlerFunc

	// 转换器
	// 为空则使用默认转换函数
	Convertor func(HandlerFunc) gin.HandlerFunc

	// 子接口组
	Groups []Group
}

// 设置接口
func Handle(r gin.IRouter, method, path string, handler gin.HandlerFunc) {
	if method == "ANY" {
		r.Any(path, handler)
	} else {
		r.Handle(method, path, handler)
	}
}

// 绑定接口
func (group Group) Bind(r gin.IRouter) {
	if len(group.Middlewares) != 0 {
		r.Use(group.Middlewares...)
	}
	if group.CustomFunc != nil {
		group.CustomFunc(r)
	}
	for _, handler := range group.Handlers {
		method, path := SplitName(handler)
		if group.Convertor != nil {
			Handle(r, method, path, group.Convertor(handler))
		} else {
			Handle(r, method, path, DefaultConvertor(handler))
		}
	}
	for path, handler := range group.HandlerMap {
		method, _ := SplitName(handler)
		if group.Convertor != nil {
			Handle(r, method, path, group.Convertor(handler))
		} else {
			Handle(r, method, path, DefaultConvertor(handler))
		}
	}
	for _, v := range group.Groups {
		v.Bind(r.Group(v.Path))
	}
}
