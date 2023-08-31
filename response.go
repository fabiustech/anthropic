package anthropic

// Response represents the response from the API.
type Response struct {
	// Completion is he resulting completion up to and excluding the stop sequences.
	Completion string `json:"completion"`
	// StopReason is the reason Anthropic stopped sampling. It will be one of "stop_sequence" or "max_tokens".
	StopReason *string `json:"stop_reason"`
	// Stop is the stop sequence that caused the model to stop sampling.
	Stop *string `json:"stop"`
	// Model is the model that performed the completion.
	Model Model `json:"model"`
}
