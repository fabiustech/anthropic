package anthropic

import (
	"fmt"
	"net/http"
)

const (
	errRateLimit  = "rate_limit_error"
	errOverloaded = "overloaded_error"
)

type ResponseError struct {
	Err Error `json:"error"`
}

// Error implements the error interface.
func (r *ResponseError) Error() string {
	return fmt.Sprintf("%s: %s (code: %d)", r.Err.Type, r.Err.Message, r.Err.Code)
}

// Retryable returns true if the error is retryable. For now, we assume all 5xx errors are transient.
// We also return true on 429, but it's up to the caller to determine the appropriate retry strategy
// given rate limits are based on # of concurrent requests.
func (r *ResponseError) Retryable() bool {
	return r.Err.Code >= http.StatusInternalServerError ||
		r.Err.Code == http.StatusTooManyRequests ||
		r.Err.Type == errOverloaded ||
		r.Err.Type == errRateLimit
}

// Error represents an error returned from the API.
type Error struct {
	// Type represents the type of error (e.g. "invalid_request_error").
	Type string `json:"type"`
	// Message is a human-readable message about the error.
	Message string `json:"message"`
	// Code is the HTTP status code returned by the API (populated by the client).
	Code int `json:"code"`
}
