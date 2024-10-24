package group

import (
	"regexp"
	"strings"
	"sync"

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

var MethodExpr = regexp.MustCompile(`\.(` + strings.Join(MethodAny, "|") + `)(\w+)`)

//go:linkname NameOfFunction github.com/gin-gonic/gin.nameOfFunction
func NameOfFunction(any) string

var relativePath *strings.Replacer

func init() {
	oldnew := []string{"ID", "_id"}
	for i := 'A'; i <= 'Z'; i++ {
		oldnew = append(oldnew, string(i)+"ID", "_"+string(i+32)+"id", string(i), "_"+string(i+32))
	}
	relativePath = strings.NewReplacer(oldnew...)
}

func ParsePath(path string) string {
	return relativePath.Replace(path)[1:]
}

var pathCache sync.Map // map[string][2]string

func Wrapper(method, path string, handler HandlerFunc) HandlerFunc {
	pathCache.Store(NameOfFunction(handler), [2]string{method, path})
	return handler
}
