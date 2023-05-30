package anthropic

// Request represents the request to the API.
type Request struct {
	// Prompt is the prompt you want Claude to complete. Required.
	Prompt Prompt `json:"prompt"`
	// Model controls which version of Claude answers your request. For more on the models, see the documentation in
	// models.go or visit https://console.anthropic.com/docs/api/reference. Required.
	Model Model `json:"model"`
	// MaxTokensToSample is the maximum number of tokens to generate before stopping. Required.
	MaxTokensToSample int `json:"max_tokens_to_sample"`
	// StopSequences specifies a list of sequences to stop sampling at. Anthropic's models stop on "\n\nHuman:", and
	// may include additional built-in stop sequences in the future. By providing the stop_sequences parameter, you may
	// include additional strings that will cause the model to stop generating.
	StopSequences []string `json:"stop_sequences,omitempty"`
	// Temperature specifies the amount of randomness injected into the response. Ranges from 0 to 1. Use temp closer to
	// 0 for analytical / multiple choice, and temp closer to 1 for creative and generative tasks.
	// Optional. Defaults to 1.
	Temperature *float64 `json:"temperature,omitempty"`
	// TopK specifies to only sample from the top K options for each subsequent token. Used to remove "long tail" low
	// probability responses.
	// Optional. Defaults to -1, which disables it. You should either alter Temperature or TopP, but not both.
	TopK *int `json:"topK,omitempty"`
	// TopP does nucleus sampling, in which we compute the cumulative distribution over all the options for each
	// subsequent token in decreasing probability order and cut it off once it reaches a particular probability
	// specified by TopP.
	//	Optional: Defaults to -1, which disables it. You should either alter Temperature or TopP, but not both.
	TopP *int `json:"topP,omitempty"`
	// Metadata is an object describing metadata about the request. Optional.
	Metadata *Metadata `json:"metadata,omitempty"`
}

// Metadata is an object describing metadata about the request.
type Metadata struct {
	// UserID is a UUID, hash value, or other external identifier for the user who is associated with the request.
	// Anthropic may use this id to help detect abuse. Do not include any identifying information such as name, email
	// address, or phone number.
	UserID string `json:"user_id,omitempty"`
}
