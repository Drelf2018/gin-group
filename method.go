package group

import (
	"regexp"
	"strings"

	_ "unsafe"
)

const (
	MethodGet     = "Get"
	MethodHead    = "Head"
	MethodPost    = "Post"
	MethodPut     = "Put"
	MethodPatch   = "Patch" // RFC 5789
	MethodDelete  = "Delete"
	MethodConnect = "Connect"
	MethodOptions = "Options"
	MethodTrace   = "Trace"
)

var MethodAny = []string{
	MethodGet,
	MethodHead,
	MethodPost,
	MethodPut,
	MethodPatch,
	MethodDelete,
	MethodConnect,
	MethodOptions,
	MethodTrace,
}

// 请求方法的正则表达式
var MethodExpr = regexp.MustCompile(`\.(Any|` + strings.Join(MethodAny, "|") + `)(\w*)`)

// 获取函数名
//
//go:linkname NameOfFunction github.com/gin-gonic/gin.nameOfFunction
func NameOfFunction(any) string

//go:linkname relativePath
var relativePath *strings.Replacer

func init() {
	oldnew := []string{"ID", "/:id"}
	for i := 'A'; i <= 'Z'; i++ {
		oldnew = append(oldnew, string(i)+"ID", "/:"+string(i+32)+"id", string(i), "/:"+string(i+32))
	}
	relativePath = strings.NewReplacer(oldnew...)
}

// 路径解析
var ParsePath = func(path string) string {
	if path == "" {
		return ""
	}
	new := relativePath.Replace(path)
	if 'A' <= path[0] && path[0] <= 'Z' {
		new = new[2:]
	}
	return "/" + new
}

// 分割接口名
func SplitName(handler HandlerFunc) (method, path string) {
	name := NameOfFunction(handler)
	matched := MethodExpr.FindStringSubmatch(name)
	if len(matched) == 3 && matched[1] != "" {
		method = strings.ToUpper(matched[1])
		path = ParsePath(matched[2])
	}
	return
}
