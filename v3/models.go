package v3

// Model represents all models.
type Model int

const (
	// UnknownModel represents an invalid model.
	UnknownModel Model = iota

	// Claude3Opus20240229 is Anthropic's most powerful model, delivering state-of-the-art performance on highly complex
	// tasks and demonstrating fluency and human-like understanding.
	Claude3Opus20240229

	// Claude3Sonnet20240229 is Anthropic's most balanced model between intelligence and speed, a great choice for
	// enterprise workloads and scaled AI deployments.
	Claude3Sonnet20240229

	// Claude3Haiku20240307 is Anthropic's fastest and most compact model, designed for near-instant responsiveness and
	// seamless AI experiences that mimic human interactions.
	Claude3Haiku20240307

	// Claude3Dot5Sonnet20240620 is Anthropic's 3.5 model that is supposed to be better than Opus while being the same speed and price as Sonnet.
	Claude3Dot5Sonnet20240620

	// Claude3Dot5Sonnet20241022 This version shows significant improvements in coding capabilities, improving performance on SWE-bench Verified from 33.4% to 49.0%, scoring higher than all publicly available models 2. The upgraded Claude 3.5 Sonnet delivers these improvements while maintaining the same price and speed as its predecessor .
	Claude3Dot5Sonnet20241022

	// Claude3Dot7Sonnet20250219 is Anthropics most intelligent model, which allows for extended thinking.
	Claude3Dot7Sonnet20250219
)

// String implements the fmt.Stringer interface.
func (c Model) String() string {
	return completionToString[c]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c Model) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// On unrecognized value, it sets |e| to Unknown.
func (c *Model) UnmarshalText(b []byte) error {
	if val, ok := stringToCompletion[(string(b))]; ok {
		*c = val
		return nil
	}

	*c = UnknownModel

	return nil
}

var completionToString = map[Model]string{
	Claude3Opus20240229:       "claude-3-opus-20240229",
	Claude3Sonnet20240229:     "claude-3-sonnet-20240229",
	Claude3Haiku20240307:      "claude-3-haiku-20240307",
	Claude3Dot5Sonnet20240620: "claude-3-5-sonnet-20240620",
	Claude3Dot5Sonnet20241022: "claude-3-5-sonnet-20241022",
}

var stringToCompletion = map[string]Model{
	"claude-3-opus-20240229":     Claude3Opus20240229,
	"claude-3-sonnet-20240229":   Claude3Sonnet20240229,
	"claude-3-haiku-20240307":    Claude3Haiku20240307,
	"claude-3-5-sonnet-20240620": Claude3Dot5Sonnet20240620,
	"claude-3-5-sonnet-20241022": Claude3Dot5Sonnet20241022,
}
