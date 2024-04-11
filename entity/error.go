package entity

type ErrorCode interface {
	error
	ErrCode() int
}

type DefaultErrorCode struct {
	Code int    `json:"code,omitempty"`
	Info string `json:"info,omitempty"`
}

func (e *DefaultErrorCode) ErrorCode() int {
	return e.Code
}

func (e *DefaultErrorCode) Error() string {
	return e.Info
}

type Response struct {
	Error *DefaultErrorCode `json:"error,omitempty"`
	Data  interface{}       `json:"data"`
}

var ErrSuccess *DefaultErrorCode = &DefaultErrorCode{Code: 0, Info: ""}
var ErrUnknown *DefaultErrorCode = &DefaultErrorCode{Code: 10000, Info: "unknown error"}
var ErrInvalidArguments *DefaultErrorCode = &DefaultErrorCode{Code: 10001, Info: "invalid arguments"}
var ErrSessionNotFound *DefaultErrorCode = &DefaultErrorCode{Code: 10002, Info: "shell not found"}
var ErrUnauthorized *DefaultErrorCode = &DefaultErrorCode{Code: 10003, Info: "unauthorized"}
var ErrHandlerNotFound *DefaultErrorCode = &DefaultErrorCode{Code: 10004, Info: "handler not found"}
var ErrFileCorrupted *DefaultErrorCode = &DefaultErrorCode{Code: 10005, Info: "file corrupted"}
var ErrFileExisted *DefaultErrorCode = &DefaultErrorCode{Code: 10006, Info: "file existed"}
