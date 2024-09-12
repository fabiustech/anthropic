package v3

import "encoding/json"

// Message represents a message sent to the API.
type Message struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message. It can be either a string or a media source.
	Content []*MessageContent `json:"content"`
}

// ShortHandMessage represents a message sent to the API in a shorthand format.
// It only supports text.
type ShortHandMessage struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
}

// MessageContent represents the content of a message.
type MessageContent struct {
	// Type is the type of the content. It can be either "text", "image", or "tool_use", or "tool_result"
	// ("tool_result" is only used when there's an error with the tool usage by the model and the model is being instructed to fix it in a subsequent call).
	Type string `json:"type"`
	// Text is the text content of the message. Leave this empty if passing an image.
	Text string `json:"text,omitempty"`
	// Source is the media source of the message. Leave this empty if passing text.
	Source *MediaSource `json:"source,omitempty"`
	// Name is the name of the tool used (if any) .
	Name string `json:"name,omitempty"`
	// Input is the input of for a specified tool (if any).
	Input json.RawMessage `json:"input,omitempty"`
	// Content is the result of a calling specified tool (if any).
	Content string `json:"content,omitempty"`
	// IsError is true only when there is an error with the first tool usage and the model is being instructed to try again.
	IsError bool `json:"is_error,omitempty"`
	// ToolUseID is the ID of the tool usage, only used when the model is instructed to try again.
	ToolUseID string `json:"tool_use_id,omitempty"`
}

// MediaSource represents the media source of a message.
type MediaSource struct {
	// Type is the type of the media source.
	Type string `json:"type"`
	// MediaType is the media type of the media source.
	MediaType string `json:"media_type"`
	// Data is the data of the media source.
	Data string `json:"data"`
}

// CacheControl represents the cache control of a system message.
type CacheControl struct {
	// Type is the type of the cache control. Currently only "ephemeral" is supported.
	Type string `json:"type"`
}

// SystemMessage represents a system message.
type SystemMessage struct {
	// Type is the type of the system message. Currently only "text" is supported.
	Type string `json:"type"`
	// Text is the text content of the system message.
	Text string `json:"text"`
	// CacheControl is the cache control of the system message. If set, the system message will not be cached (assuming
	// the prompt caching beta header is included in the request.
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}
