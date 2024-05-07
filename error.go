package group

import (
	"fmt"
	"runtime"
)

type Error interface {
	error
	Code() int
}

type ErrorWithCode struct {
	code    int
	message string
}

func (e ErrorWithCode) Code() int {
	return e.code
}

func (e ErrorWithCode) Error() string {
	return e.message
}

func (e ErrorWithCode) Format(a ...any) ErrorWithCode {
	e.message = fmt.Sprintf(e.message, a...)
	return e
}

var _ Error = ErrorWithCode{}

func NewError(code int, message string) ErrorWithCode {
	return ErrorWithCode{code, message}
}

func AutoError(err error) (e ErrorWithCode) {
	_, _, e.code, _ = runtime.Caller(1)
	e.message = err.Error()
	return
}
