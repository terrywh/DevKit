package entity

type ErrorWithCode interface {
	ErrorCode() int
	Error() string
}

type HttpError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

func (e HttpError) ErrorCode() int {
	return e.Code
}

func (e HttpError) Error() string {
	return e.Info
}

type HttpResponse struct {
	Error HttpError   `json:"error"`
	Data  interface{} `json:"data"`
}

var ErrSuccess HttpError = HttpError{Code: 0, Info: ""}
var ErrUnknown HttpError = HttpError{Code: 10000, Info: "unknown error"}
var ErrInvalidArguments HttpError = HttpError{Code: 10001, Info: "invalid arguments"}
var ErrSessionNotFound HttpError = HttpError{Code: 10002, Info: "shell not found"}
var ErrUnauthorized HttpError = HttpError{Code: 10003, Info: "unauthorized"}
var ErrHandlerNotFound HttpError = HttpError{Code: 10004, Info: "handler not found"}
var ErrFileCorrupted HttpError = HttpError{Code: 10005, Info: "file corrupted"}
