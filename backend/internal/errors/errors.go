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

// HTTPError returns an HTTP error for the middleware to send (e.g. 400 Bad Request, 500 Internal Server Error).
func HTTPError(code int, err error) *httpErrorT {
	return newError(code, err)
}

const (
	myCustomErrorCode    = 505
	myCustomErrorMessage = "my custom error"
)

func MyCustomError() *httpErrorT {
	return newError(myCustomErrorCode, errors.New(myCustomErrorMessage))
}
