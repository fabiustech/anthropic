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

	// Claude3Dot5Sonnet is the latest version of the Claude 3.5 Sonnet model.
	Claude3Dot5Sonnet

	// Claude3Dot5Opus is the latest version of the Claude 3.5 Opus model.
	Claude3Dot5Opus

	// Claude3Dot5Haiku is the latest version of the Claude 3.5 Haiku model.
	Claude3Dot5Haiku

	// Claude3Dot7Sonnet is the latest version of the Claude 3.7 Sonnet model.
	Claude3Dot7Sonnet

	// Claude4Sonnet20250514 is Anthropic's latest high-performance model with exceptional reasoning and efficiency.
	Claude4Sonnet20250514

	// Claude4Opus20250514 is Anthropic's most capable and intelligent model yet. Claude Opus 4 sets new standards in
	// complex reasoning and advanced coding.
	Claude4Opus20250514

	// ClaudeSonnet4Dot5 is Anthropic's smartest model for complex agents and coding.
	ClaudeSonnet4Dot5

	// ClaudeHaiku4Dot5 is Anthropic's fastest model with near-frontier intelligence.
	ClaudeHaiku4Dot5

	// ClaudeOpus4Dot1 is an exceptional model for specialized reasoning tasks.
	ClaudeOpus4Dot1
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
	Claude3Dot7Sonnet20250219: "claude-3-7-sonnet-20250219",
	Claude3Dot5Sonnet:         "claude-3-5-sonnet-latest",
	Claude3Dot5Opus:           "claude-3-5-opus-latest",
	Claude3Dot5Haiku:          "claude-3-5-haiku-latest",
	Claude3Dot7Sonnet:         "claude-3-7-sonnet-latest",
	Claude4Sonnet20250514:     "claude-sonnet-4-20250514",
	Claude4Opus20250514:       "claude-opus-4-20250514",
	ClaudeSonnet4Dot5:         "claude-sonnet-4-5",
	ClaudeHaiku4Dot5:          "claude-haiku-4-5",
	ClaudeOpus4Dot1:           "claude-opus-4-1",
}

var stringToCompletion = map[string]Model{
	"claude-3-opus-20240229":     Claude3Opus20240229,
	"claude-3-sonnet-20240229":   Claude3Sonnet20240229,
	"claude-3-haiku-20240307":    Claude3Haiku20240307,
	"claude-3-5-sonnet-20240620": Claude3Dot5Sonnet20240620,
	"claude-3-5-sonnet-20241022": Claude3Dot5Sonnet20241022,
	"claude-3-7-sonnet-20250219": Claude3Dot7Sonnet20250219,
	"claude-3-5-sonnet-latest":   Claude3Dot5Sonnet,
	"claude-3-5-opus-latest":     Claude3Dot5Opus,
	"claude-3-5-haiku-latest":    Claude3Dot5Haiku,
	"claude-3-7-sonnet-latest":   Claude3Dot7Sonnet,
	"claude-sonnet-4-20250514":   Claude4Sonnet20250514,
	"claude-opus-4-20250514":     Claude4Opus20250514,
	"claude-sonnet-4-5":          ClaudeSonnet4Dot5,
	"claude-haiku-4-5":           ClaudeHaiku4Dot5,
	"claude-opus-4-1":            ClaudeOpus4Dot1,
}
