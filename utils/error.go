package utils

type YYError interface {
	error
	ErrorCode()
	ErrorMsg()
}
