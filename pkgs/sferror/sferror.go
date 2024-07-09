package sferror

import (
	"errors"
	"fmt"
	"runtime"
)

// A Error 基础错误
type Error struct {
	Code    int
	Err     error
	FileNum string
}

// Error 返回异常信息
func (err *Error) Error() string {
	return err.Err.Error()
}

// New 创建一个新的error
func (err *Error) New() *Error {
	_, file, line, ok := runtime.Caller(1)
	newErr := Error{
		Code: err.Code,
		Err:  err.Err,
	}
	if ok {
		newErr.FileNum = fmt.Sprintf("%s:%d", file, line)
	}
	return &newErr
}

// NewError 创建一个新Error
func NewError(code int, err error) *Error {
	if err == nil {
		err = errors.New("system error")
	}
	return &Error{Code: code, Err: err}
}

// NewUndefinedError error
func NewUndefinedError() *Error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{Code: 1, Err: errors.New("undefined error"), FileNum: fmt.Sprintf("%s:%d", file, line)}
}

// CheckAndPanic 检查错误并panic
// 这是非常不建议使用的一种方式，会掩盖错误的正常处理流程
func CheckAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
