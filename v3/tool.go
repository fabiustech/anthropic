package v3

// Tool represents a tool that the model may use.
type Tool struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	InputSchema *Schema `json:"input_schema"`
}

// Schema represents a basic JSON schema.
type Schema struct {
	Type        SchemaType         `json:"type"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Title       string             `json:"title,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty"`
	Description string             `json:"description,omitempty"`
	Required    []string           `json:"required,omitempty"`
}

// ToolChoice represents how the model should use the provided tools. The model can use a specific tool, any available
// tool, or decide by itself.
type ToolChoice struct {
	// Type is the type of tool choice. Required.
	Type string `json:"type"`
	// Name is the name of the tool to use. Required if Type is "tool".
	Name string `json:"name,omitempty"`
}

// ToolChoiceType represents the type of tool choice.
type ToolChoiceType int

const (
	// ToolChoiceAuto allows Claude to decide whether to call any provided tools or not. This is the default value.
	ToolChoiceAuto ToolChoiceType = iota
	// ToolChoiceAny tells Claude that it must use one of the provided tools, but doesn’t force a particular tool.
	ToolChoiceAny
	// ToolChoiceTool allows us to force Claude to always use a particular tool.
	ToolChoiceTool
)

// String implements the fmt.Stringer interface.
func (t ToolChoiceType) String() string {
	return toolChoiceTypeToString[t]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t ToolChoiceType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// On unrecognized value, it sets |e| to ToolChoiceAuto (the default).
func (t *ToolChoiceType) UnmarshalText(b []byte) error {
	if val, ok := stringToToolChoiceType[(string(b))]; ok {
		*t = val
		return nil
	}

	*t = ToolChoiceAuto

	return nil
}

var toolChoiceTypeToString = map[ToolChoiceType]string{
	ToolChoiceAuto: "auto",
	ToolChoiceAny:  "any",
	ToolChoiceTool: "tool",
}

var stringToToolChoiceType = map[string]ToolChoiceType{
	"auto": ToolChoiceAuto,
	"any":  ToolChoiceAny,
	"tool": ToolChoiceTool,
}

// SchemaType represents the type of a JSON schema.
type SchemaType int

const (
	// SchemaTypeNull represents a null value.
	SchemaTypeNull SchemaType = iota
	// SchemaTypeString represents a string value.
	SchemaTypeString
	// SchemaTypeNumber represents a number value.
	SchemaTypeNumber
	// SchemaTypeInteger represents an integer value.
	SchemaTypeInteger
	// SchemaTypeObject represents an object value.
	SchemaTypeObject
	// SchemaTypeArray represents an array value.
	SchemaTypeArray
	// SchemaTypeBoolean represents a boolean value.
	SchemaTypeBoolean
)

// String implements the fmt.Stringer interface.
func (st SchemaType) String() string {
	return schemaTypeToString[st]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (st SchemaType) MarshalText() ([]byte, error) {
	return []byte(st.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// On unrecognized value, it sets |e| to Null.
func (st *SchemaType) UnmarshalText(b []byte) error {
	if val, ok := stringToSchemaType[(string(b))]; ok {
		*st = val
		return nil
	}

	*st = SchemaTypeNull

	return nil
}

var schemaTypeToString = map[SchemaType]string{
	SchemaTypeNull:    "null",
	SchemaTypeString:  "string",
	SchemaTypeNumber:  "number",
	SchemaTypeInteger: "integer",
	SchemaTypeObject:  "object",
	SchemaTypeArray:   "array",
	SchemaTypeBoolean: "boolean",
}

var stringToSchemaType = map[string]SchemaType{
	"null":    SchemaTypeNull,
	"string":  SchemaTypeString,
	"number":  SchemaTypeNumber,
	"integer": SchemaTypeInteger,
	"object":  SchemaTypeObject,
	"array":   SchemaTypeArray,
	"boolean": SchemaTypeBoolean,
}
