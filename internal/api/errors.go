package api

import "fmt"

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
	ExitUsage   = 2
	ExitAuth    = 3
	ExitNetwork = 4
	ExitNotFound = 5
)

// NetworkError represents a network connectivity error
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// APIError represents an error from the API
type APIError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// ExitCode returns the appropriate exit code for the error
func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.StatusCode {
		case 401, 403:
			return ExitAuth
		case 404:
			return ExitNotFound
		default:
			return ExitError
		}
	}

	if _, ok := err.(*NetworkError); ok {
		return ExitNetwork
	}

	return ExitError
}
