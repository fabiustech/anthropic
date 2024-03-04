package v3

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
	// Type is the type of the content. It can be either "text" or "image".
	Type string `json:"type"`
	// Text is the text content of the message. Leave this empty if passing an image.
	Text string `json:"text,omitempty"`
	// Source is the media source of the message. Leave this empty if passing text.
	Source *MediaSource `json:"source,omitempty"`
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
