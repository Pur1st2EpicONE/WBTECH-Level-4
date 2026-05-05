// Package errs defines shared application-level errors.
package errs

import "errors"

var (
	ErrInternal = errors.New("internal server error") // internal server error
)
