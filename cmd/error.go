package cmd

import (
	"fmt"
	"net/http"
	"strings"
)

// Errors is a list of errors.
// Useful in a loop if you don't want to return the error right away and you want to display after the loop,
// all the errors that happened during the loop.
type Errors []error

func (errList Errors) Error() string {
	if len(errList) < 1 {
		return ""
	}

	out := make([]string, len(errList))
	for i := range errList {
		out[i] = errList[i].Error()
	}
	return strings.Join(out, ", ")
}

// StatusError reports an unsuccessful exit by a command.
type StatusError struct {
	Status     string
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("Status: %s, Code: %d", e.Status, e.StatusCode)
}

type statusError struct {
	status int
	method string
	err    string
}

// Error returns a formatted string for this error type
func (e *statusError) Error() string {
	return fmt.Sprintf("%s: %v", e.method, e.err)
}

// IsNotFound indicates if the passed in error is from an http.StatusNotFound from the plugin
func IsNotFound(err error) bool {
	return isStatusError(err, http.StatusNotFound)
}

func isStatusError(err error, status int) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*statusError)
	if !ok {
		return false
	}
	return e.status == status
}
