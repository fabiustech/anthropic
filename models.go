package anthropic

// Model represents all models.
type Model int

const (
	// UnknownModel represents an invalid model.
	UnknownModel Model = iota
	// Claude is Anthropic's largest model, ideal for a wide range of more complex tasks. This is the "major version"
	// which will automatically get updates to the model as they are released.
	Claude
	// Claude2Dot0 is Anthropic's largest model, ideal for a wide range of more complex tasks. If you rely on the exact
	// output shape, you should specify this full model version. It has a context window of 100K tokens.
	Claude2Dot0
	// Claude2Dot1 represents an improvement in specific capabilities and performance over Claude 2. With strong
	// accuracy upgrades, double the context window, and experimental tool use features, Claude can handle more complex
	// reasoning and generation while remaining honest and grounded in factual evidence.
	// Claude 2.1's context window is 200K tokens, enabling it to leverage much richer contextual information to
	// generate higher quality and more nuanced output.
	Claude2Dot1
	// ClaudeInstant is a smaller model with far lower latency, sampling at roughly 40 words/sec! Its output quality
	// is somewhat lower than the latest Claude model, particularly for complex tasks. However, it is much less
	// expensive and blazing fast. Anthropic believes that this model provides more than adequate performance on a range
	// of tasks including text classification, summarization, and lightweight chat applications, as well as search
	// result summarization. This is the "major version" which will automatically get updates to the model as they
	// are released.
	ClaudeInstant
	// ClaudeInstant1Dot1 is a smaller model with far lower latency, sampling at roughly 40 words/sec! Its output
	// quality is somewhat lower than the latest Claude model, particularly for complex tasks. However, it is much less
	// expensive and blazing fast. Anthropic believes that this model provides more than adequate performance on a range
	// of tasks including text classification, summarization, and lightweight chat applications, as well as search
	// result summarization. If you rely on the exact output shape, you should specify this full model version.
	ClaudeInstant1Dot1

	// Claude3Opus20240229 is Anthropic's most powerful model, delivering state-of-the-art performance on highly complex
	// tasks and demonstrating fluency and human-like understanding.
	Claude3Opus20240229

	// Claude3Sonnet20240229 is Anthropic's most balanced model between intelligence and speed, a great choice for
	// enterprise workloads and scaled AI deployments.
	Claude3Sonnet20240229
)

// String implements the fmt.Stringer interface.
func (c Model) String() string {
	return completionToString[c]
}

// BedrockString returns the string representation of the model for use with AWS Bedrock.
func (c Model) BedrockString() string {
	return bedrockToString[c]
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
	Claude:                "claude-2",
	Claude2Dot0:           "claude-2.0",
	Claude2Dot1:           "claude-2.1",
	ClaudeInstant:         "claude-instant-1",
	ClaudeInstant1Dot1:    "claude-instant-1.1",
	Claude3Opus20240229:   "claude-3-opus-20240229",
	Claude3Sonnet20240229: "claude-3-sonnet-20240229",
}

var stringToCompletion = map[string]Model{
	"claude-2":                 Claude,
	"claude-2.0":               Claude2Dot0,
	"claude-2.1":               Claude2Dot1,
	"claude-instant-1":         ClaudeInstant,
	"claude-instant-1.1":       ClaudeInstant1Dot1,
	"claude-3-opus-20240229":   Claude3Opus20240229,
	"claude-3-sonnet-20240229": Claude3Sonnet20240229,
}

var bedrockToString = map[Model]string{
	Claude:        "anthropic.claude-v2",
	Claude2Dot0:   "anthropic.claude-v2",
	ClaudeInstant: "anthropic.claude-instant-v1",
}
