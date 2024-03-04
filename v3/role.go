package v3

// Role represents the role of the message sender.
type Role int

const (
	// RoleUnknown represents an unknown role.
	RoleUnknown Role = iota
	// RoleUser represents a user role.
	RoleUser
	// RoleAssistant represents an assistant role.
	RoleAssistant
)

// String returns the string representation of the role.
func (r Role) String() string {
	return roleToString[r]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// On unrecognized value, it sets |e| to Unknown.
func (r *Role) UnmarshalText(b []byte) error {
	if val, ok := stringToRole[(string(b))]; ok {
		*r = val
		return nil
	}

	*r = RoleUnknown

	return nil
}

var roleToString = map[Role]string{
	RoleUser:      "user",
	RoleAssistant: "assistant",
}

var stringToRole = map[string]Role{
	"user":      RoleUser,
	"assistant": RoleAssistant,
}
