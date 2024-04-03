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

var ErrUnknown HttpError = HttpError{Code: 10000, Info: "unknown error"}
var ErrShellNotFound HttpError = HttpError{Code: 10001, Info: "shell not found"}
