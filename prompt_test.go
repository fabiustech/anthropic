package anthropic

import (
	"testing"
)

func TestMessageMarshal(t *testing.T) {
	var tcs = []struct {
		msg Message
		exp string
	}{
		{msg: Message{UserType: UserTypeHuman, Text: "Hello"}, exp: "\n\nHuman: Hello"},
		{msg: Message{UserType: UserTypeAssistant, Text: "How can I help?"}, exp: "\n\nAssistant: How can I help?"},
		{msg: Message{UserType: UserTypeSystem, Text: "System message"}, exp: "System message"},
	}

	for _, tc := range tcs {
		if result := tc.msg.marshal(); result != tc.exp {
			t.Errorf("marshal() = %v, want %v", result, tc.exp)
		}
	}
}

func TestMessagesValidate(t *testing.T) {
	tests := []struct {
		name string
		msgs Messages
		err  error
	}{
		{
			name: "Empty Messages",
			msgs: Messages{},
			err:  ErrEmptyMessages,
		},
		{
			name: "System Message Not First",
			msgs: Messages{
				{UserType: UserTypeHuman, Text: "Hi"},
				{UserType: UserTypeSystem, Text: "System message"},
			},
			err: ErrBadSystemMessage,
		},
		{
			name: "Missing Assistant Message",
			msgs: Messages{
				{UserType: UserTypeHuman, Text: "Hello"},
			},
			err: ErrMissingAssistant,
		},
		{
			name: "Valid Messages",
			msgs: Messages{
				{UserType: UserTypeSystem, Text: "System starting"},
				{UserType: UserTypeHuman, Text: "Hi"},
				{UserType: UserTypeAssistant, Text: "Hello"},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err = tt.msgs.Validate()
			if err != tt.err {
				t.Errorf("Messages.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestNewPromptFromMessages(t *testing.T) {
	var tcs = []struct {
		msgs []*Message
		exp  string
	}{
		{
			msgs: []*Message{
				{UserType: UserTypeHuman, Text: "Hi"},
				{UserType: UserTypeAssistant, Text: "Hello!"},
			},
			exp: "\n\nHuman: Hi\n\nAssistant: Hello!",
		},
		{
			msgs: []*Message{
				{UserType: UserTypeHuman, Text: "Hi"},
				{UserType: UserTypeAssistant, Text: "Hello!"},
				{UserType: UserTypeHuman, Text: "How are you?"},
				{UserType: UserTypeAssistant, Text: "I'm good, thanks!"},
			},
			exp: "\n\nHuman: Hi\n\nAssistant: Hello!\n\nHuman: How are you?\n\nAssistant: I'm good, thanks!",
		},
		{
			msgs: []*Message{
				{UserType: UserTypeSystem, Text: "Test system message"},
				{UserType: UserTypeHuman, Text: "Hello!"},
				{UserType: UserTypeAssistant, Text: "How can I help?"},
			},
			exp: "Test system message\n\nHuman: Hello!\n\nAssistant: How can I help?",
		},
	}

	for _, tc := range tcs {
		if result := NewPromptFromMessages(tc.msgs); string(result) != tc.exp {
			t.Errorf("NewPromptFromMessages() = %v, want %v", result, tc.exp)
		}
	}
}

func TestNewPromptFromString(t *testing.T) {
	var input = "What's the weather like?"
	var exp = "\n\nHuman: What's the weather like?\n\nAssistant:"
	if result := NewPromptFromString(input); string(result) != exp {
		t.Errorf("NewPromptFromString() = %v, want %v", result, exp)
	}
}

func TestNewPromptFromStringWithSystemMessage(t *testing.T) {
	var system = "You are a test."
	var human = "Confirm this test passes."
	var exp = "You are a test.\n\nHuman: Confirm this test passes.\n\nAssistant:"
	if result := NewPromptFromStringWithSystemMessage(system, human); string(result) != exp {
		t.Errorf("NewPromptFromStringWithSystemMessage() = %v, want %v", result, exp)
	}
}
