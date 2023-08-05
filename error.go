package anthropic

import "fmt"

type errorResponse struct {
	Err Error `json:"error"`
}

// Error implements the error interface.
func (e *errorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Err.Type, e.Err.Message)
}

// Error represents an error returned from the API.
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
