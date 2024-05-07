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

var AnyMethods = []string{
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

var MethodEXP = regexp.MustCompile("^(" + strings.Join(AnyMethods, "|") + ")(\\w+)")

//go:linkname NameOfFunction github.com/gin-gonic/gin.nameOfFunction
func NameOfFunction(any) string

func NameOfHandler(fn HandlerFunc) string {
	s := strings.Split(NameOfFunction(fn), ".")
	return s[len(s)-1]
}

func SplitHandlerName(fn HandlerFunc) (method, path string) {
	s := MethodEXP.FindStringSubmatch(NameOfHandler(fn))
	if len(s) == 3 {
		method = s[1]
		path = s[2]
	}
	return
}

const (
	Colon    = "Colon"
	Asterisk = "Asterisk"
)

var relativePath *strings.Replacer

func init() {
	oldnew := []string{Colon, "/:", Asterisk, "/*"}
	for i := 'A'; i <= 'Z'; i++ {
		oldnew = append(oldnew, "_"+string(i), string(i), string(i), "/"+string(i+32))
	}
	relativePath = strings.NewReplacer(oldnew...)
}

func ParsePath(path string) string {
	return relativePath.Replace(path)
}
