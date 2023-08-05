package anthropic

import "fmt"

type errorResponse struct {
	Err Error `json:"error"`
}

func (e *errorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Err.Type, e.Err.Message)
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
