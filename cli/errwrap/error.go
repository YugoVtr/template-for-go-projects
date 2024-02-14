// Package errors provides custom error types and error handling functions for the application.
// based on;
//   - https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
//   - https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
//   - https://dev.to/jhall/error-handling-in-go-web-apps-shouldnt-be-so-awkward-4e2k
package errwrap

import "fmt"

// Error represents a custom error type in the application.
type Error struct {
	Err         error
	Description string
}

// Error returns the string representation of the error.
// It formats the error with the description and the underlying error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Description, e.Err)
}

// Unwrap returns the underlying error wrapped by the current error.
// If the current error does not wrap another error, it returns nil.
func (e *Error) Unwrap() error {
	return e.Err
}

// WithMessage wraps an error with a description.
// It takes a description string and an error as input,
// and returns a new error with the provided description.
func WithMessage(description string, err error) error {
	return &Error{
		Err:         err,
		Description: description,
	}
}
