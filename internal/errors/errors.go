package errors

import (
	"errors"
	"fmt"
)

type httpErrorT struct {
	Code int
	Err  error
}

// Ensure httpErrorT implements error interface
var _ error = &httpErrorT{}

func (e *httpErrorT) Error() string {
	return fmt.Sprintf("code: %d, error: %v", e.Code, e.Err)
}

func (e *httpErrorT) Unwrap() error {
	return e.Err
}

func newError(code int, err error) *httpErrorT {
	return &httpErrorT{Code: code, Err: err}
}

const (
	myCustomErrorCode    = 505
	myCustomErrorMessage = "my custom error"
)

func MyCustomError() *httpErrorT {
	return newError(myCustomErrorCode, errors.New(myCustomErrorMessage))
}
