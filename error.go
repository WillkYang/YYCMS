package YYCMS

import (
	"bytes"
	"fmt"
)

type ErrorType string

const (
	NO_ERROR ErrorType = "NO_ERROR"
	DB_ERROR ErrorType = "DB_ERROR"
)

type YYError interface {
	Type() ErrorType
	Error() string
}

type YYCMSError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func NewYYError(errType ErrorType, errMsg string) YYError {
	return &YYCMSError{errType:errType, errMsg: errMsg}
}

func (yye *YYCMSError) Type() ErrorType {
	return yye.errType
}

func (yye *YYCMSError) Error() string {
	if yye.fullErrMsg == "" {
		yye.genFullErrMsg()
	}
	return yye.fullErrMsg
}

func (yye *YYCMSError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("YYError: ")
	if yye.errType != "" {
		buffer.WriteString(string(yye.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(yye.errMsg)
	yye.fullErrMsg = fmt.Sprintf("%s\n", buffer.String())
	return
}

