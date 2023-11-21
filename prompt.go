package anthropic

import (
	"fmt"
	"strings"
)

// Prompt represents the prompt passed to the model. The text that you give Claude is designed to elicit, or "prompt",
// a relevant response. A prompt is usually in the form of a question or instruction. For example:
//
//	Human: Why is the sky blue?
//
//	Assistant:
//
// Prompts sent via the API must contain \n\nHuman: and \n\nAssistant: as the signals of who's speaking. In Slack and
// our web interface we automatically add these for you.
type Prompt string

// UserType represents the two types of users in the dialogue: Human and Assistant.
type UserType string

const (
	// UserTypeHuman is the user type for the human in the dialogue.
	UserTypeHuman = "\n\nHuman"
	// UserTypeAssistant is the user type for the assistant in the dialogue.
	UserTypeAssistant = "\n\nAssistant"
	// UserTypeSystem is the user type for the system in the dialogue. It should always be the first "message" in the
	// dialogue.
	UserTypeSystem = "System"
)

// Message represents a single message in a dialogue. It contains the UserType and the text of the message.
type Message struct {
	UserType UserType
	Text     string
}

func (m *Message) marshal() string {
	if m.UserType == UserTypeSystem {
		return m.Text
	}
	return fmt.Sprintf("%s: %s", m.UserType, m.Text)
}

// NewPromptFromMessages returns a Prompt from a slice of |Message|s by wrapping them in the expected Human/Assistant
// format. You can use this style to "Put words in Claude's mouth."
// https://console.anthropic.com/docs/prompt-design#-putting-words-in-claude-s-mouth-
func NewPromptFromMessages(msg []*Message) Prompt {
	var prompt = make([]string, len(msg))
	for i, m := range msg {
		prompt[i] = m.marshal()
	}
	return Prompt(strings.Join(prompt, ""))
}

// NewPromptFromString returns a Prompt from a string by wrapping it in the expected Human/Assistant format.
func NewPromptFromString(s string) Prompt {
	return Prompt(fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", s))
}

// NewPromptFromStringWithSystemMessage returns a Prompt from both a system and human string by wrapping them in the
// exp Human/Assistant format.
func NewPromptFromStringWithSystemMessage(system, human string) Prompt {
	return Prompt(fmt.Sprintf("%s%s", system, NewPromptFromString(human)))
}
