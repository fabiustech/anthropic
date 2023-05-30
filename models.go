// Package models contains the enum values which represent the various
// models available via the Anthropic API.
package anthropic

// Model represents all models.
type Model int

const (
	// UnknownModel represents an invalid model.
	UnknownModel Model = iota
	// ClaudeV1 is Anthropic's largest model, ideal for a wide range of more complex tasks.
	ClaudeV1
	// ClaudeV1_100K is an enhanced version of claude-v1 with a 100,000 token (roughly 75,000 word) context window. Ideal
	// for summarizing, analyzing, and querying long documents and conversations for nuanced understanding of complex
	// topics and relationships across very long spans of text.
	ClaudeV1_100K
	// ClaudeInstantV1 is a smaller model with far lower latency, sampling at roughly 40 words/sec! Its output quality
	// is somewhat lower than the latest claude-v1 model, particularly for complex tasks. However, it is much less
	// expensive and blazing fast. We believe that this model provides more than adequate performance on a range of
	// tasks including text classification, summarization, and lightweight chat applications, as well as search result
	// summarization.
	ClaudeInstantV1
	// ClaudeInstantV1_100K is n enhanced version of claude-instant-v1 with a 100,000 token context window that retains
	// its performance. Well-suited for high throughput use cases needing both speed and additional context, allowing
	// deeper understanding from extended conversations and documents.
	ClaudeInstantV1_100K
	// ClaudeV1_3 is claude-v1.3. ompared to claude-v1.2, it's more robust against red-team inputs, better at precise
	// instruction-following, better at code, and better and non-English dialogue and writing.
	ClaudeV1_3
	// ClaudeV1_3_100K is n enhanced version of claude-v1.3 with a 100,000 token (roughly 75,000 word) context window.
	ClaudeV1_3_100K
	// ClaudeV1_2 is an improved version of claude-v1. It is slightly improved at general helpfulness, instruction
	// following, coding, and other tasks. It is also considerably better with non-English languages. This model also
	// has the ability to role play (in harmless ways) more consistently, and it defaults to writing somewhat longer
	// and more thorough responses.
	ClaudeV1_2
	// ClaudeV1_0 is an earlier version of claude-v1.
	ClaudeV1_0
	// ClaudeInstantV1_1 is Anthropic's latest version of claude-instant-v1. It is better than claude-instant-v1.0 at a
	// wide variety of tasks including writing, coding, and instruction following. It performs better on academic
	// benchmarks, including math, reading comprehension, and coding tests. It is also more robust against red-teaming
	// inputs.
	ClaudeInstantV1_1
	// ClaudeInstantV1_1_100K is an enhanced version of claude-instant-v1.1 with a 100,000 token context window that
	// retains its lightning fast 40 word/sec performance.
	ClaudeInstantV1_1_100K
	// ClaudeInstantV1_0 is an earlier version of claude-instant-v1.
	ClaudeInstantV1_0
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
	ClaudeV1:               "claude-v1",
	ClaudeV1_100K:          "claude-v1-100k",
	ClaudeInstantV1:        "claude-instant-v1",
	ClaudeInstantV1_100K:   "claude-instant-v1-100k",
	ClaudeV1_3:             "claude-v1.3",
	ClaudeV1_3_100K:        "claude-v1.3-100k",
	ClaudeV1_2:             "claude-v1.2",
	ClaudeV1_0:             "claude-v1.0",
	ClaudeInstantV1_1:      "claude-instant-v1.1",
	ClaudeInstantV1_1_100K: "claude-instant-v1.1-100k",
	ClaudeInstantV1_0:      "claude-instant-v1.0",
}

var stringToCompletion = map[string]Model{
	"claude-v1":                ClaudeV1,
	"claude-v1-100k":           ClaudeV1_100K,
	"claude-instant-v1":        ClaudeInstantV1,
	"claude-instant-v1-100k":   ClaudeInstantV1_100K,
	"claude-v1.3":              ClaudeV1_3,
	"claude-v1.3-100k":         ClaudeV1_3_100K,
	"claude-v1.2":              ClaudeV1_2,
	"claude-v1.0":              ClaudeV1_0,
	"claude-instant-v1.1":      ClaudeInstantV1_1,
	"claude-instant-v1.1-100k": ClaudeInstantV1_1_100K,
	"claude-instant-v1.0":      ClaudeInstantV1_0,
}
