package errors

import "fmt"

type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e Error) Status() int {
	return e.StatusCode
}

func EchoToHTTPError(code int, message interface{}) Error {
	return Error{
		Code:       fmt.Sprintf("%d", code),
		Message:    fmt.Sprintf("%v", message),
		StatusCode: code,
	}
}
