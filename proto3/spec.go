package proto3

import (
	"bytes"
	"errors"
	"fmt"
)

type ImportType string
type NameType string
type TagType uint8
type FieldType uint8

const (
	DOUBLE_TYPE FieldType = iota + 1
	FLOAT_TYPE
	INT32_TYPE
	INT64_TYPE
	UINT32_TYPE
	UINT64_TYPE
	SINT32_TYPE
	SINT64_TYPE
	FIXED32_TYPE
	FIXED64_TYPE
	SFIXED32_TYPE
	SFIXED64_TYPE
	BOOL_TYPE
	STRING_TYPE
	BYTES_TYPE
)

type Validated interface {
	Validate() error
}

type Reserved interface {
	Write() (string, error)
}

type Field interface {
	Write() (string, error)
}

type Spec struct {
	Package  string
	Imports  []ImportType
	Messages []Message
}

type Message struct {
	Name           string
	Messages       []Message
	ReservedValues []Reserved
	Fields         []Field
	Enums          []Enum
}

type ReservedName struct {
	Name NameType
}

type ReservedTagValue struct {
	Tag TagType
}

type ReservedTagRange struct {
	LowerTag TagType
	UpperTag TagType
}

type CustomField struct {
	Name     NameType
	Tag      TagType
	Repeated bool
	Comment  string
	Typing   string
}

type ScalarField struct {
	Name     NameType
	Tag      TagType
	Repeated bool
	Comment  string
	Typing   FieldType
}

type MapField struct {
	Name        NameType
	Tag         TagType
	Comment     string
	KeyTyping   FieldType
	ValueTyping FieldType
}

type CustomMapField struct {
	Name        NameType
	Tag         TagType
	Comment     string
	KeyTyping   FieldType
	ValueTyping string
}

type Enum struct {
	Name       NameType
	Values     []EnumValue
	AllowAlias bool
}

type EnumValue struct {
	Name    NameType
	Tag     TagType
	Comment string
}

// WRITERS

func (s *Spec) Write(level int) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("syntax = \"proto3\";\n")
	if s.Package != "" {
		buffer.WriteString(fmt.Sprintf("package %s;\n", s.Package))
	}
	for _, importPackage := range s.Imports {
		buffer.WriteString(fmt.Sprintf("import \"%s\";\n", importPackage))
	}
	for _, msg := range s.Messages {
		msgSpec, err := msg.Write(level)
		if err != nil {
			return "", err
		}
		buffer.WriteString(fmt.Sprintf("%s\n", msgSpec))
	}
	return buffer.String(), nil
}

func (m *Message) Write(level int) (string, error) {
	err := m.Validate()
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%smessage %s {\n", indentLevel(level), m.Name))

	// NESTED MESSAGE TYPES
	for _, msg := range m.Messages {
		msgSpec, err := msg.Write(level + 1)
		if err != nil {
			return "", err
		}
		buffer.WriteString(fmt.Sprintf("%s\n", msgSpec))
	}

	// ENUMS
	if len(m.Enums) > 0 {
		buffer.WriteString("\n")
	}
	for _, v := range m.Enums {
		v, err := v.Write(level + 1)
		if err != nil {
			return "", err
		}
		buffer.WriteString(fmt.Sprintf("%s\n", v))
	}

	// RESERVED TAGS
	if len(m.ReservedValues) > 0 {
		buffer.WriteString("\n")
	}
	for _, reservedValue := range m.ReservedValues {
		v, err := reservedValue.Write()
		if err != nil {
			return "", err
		}
		buffer.WriteString(fmt.Sprintf("%sreserved %s;\n", indentLevel(level+1), v))
	}

	// FIELDS
	if len(m.Fields) > 0 {
		buffer.WriteString("\n")
	}
	for _, v := range m.Fields {
		v, err := v.Write()
		if err != nil {
			return "", err
		}
		buffer.WriteString(fmt.Sprintf("%s%s\n", indentLevel(level+1), v))
	}

	buffer.WriteString(fmt.Sprintf("%s}\n", indentLevel(level)))
	return buffer.String(), nil
}

func (r ReservedName) Write() (string, error) {
	return fmt.Sprintf("\"%s\"", r.Name), nil
}

func (r ReservedTagValue) Write() (string, error) {
	return fmt.Sprintf("%d", r.Tag), nil
}

func (r ReservedTagRange) Write() (string, error) {
	return fmt.Sprintf("%d to %d", r.LowerTag, r.UpperTag), nil
}

func (c CustomField) Write() (string, error) {
	v := fmt.Sprintf("%s %s = %d;", c.Typing, c.Name, c.Tag)

	if c.Repeated {
		v = fmt.Sprintf("repeated %s", v)
	}
	if c.Comment != "" {
		v = fmt.Sprintf("%s   // %s", v, c.Comment)
	}
	return v, nil
}

func (s ScalarField) Write() (string, error) {
	typeString, err := toProtobufType(s.Typing)
	if err != nil {
		return "", err
	}

	v := fmt.Sprintf("%s %s = %d;", typeString, s.Name, s.Tag)
	if s.Repeated {
		v = fmt.Sprintf("repeated %s", v)
	}
	if s.Comment != "" {
		v = fmt.Sprintf("%s   // %s", v, s.Comment)
	}
	return v, nil
}

func (m MapField) Write() (string, error) {
	var keyTypeString, valueTypeString string
	var err error
	keyTypeString, err = toProtobufType(m.KeyTyping)
	if err != nil {
		return "", err
	}
	valueTypeString, err = toProtobufType(m.ValueTyping)
	if err != nil {
		return "", err
	}
	v := fmt.Sprintf("map<%s, %s> %s = %d;", keyTypeString, valueTypeString, m.Name, m.Tag)
	if m.Comment != "" {
		v = fmt.Sprintf("%s   // %s", v, m.Comment)
	}
	return v, nil
}

func (c CustomMapField) Write() (string, error) {
	keyTypeString, err := toProtobufType(c.KeyTyping)
	if err != nil {
		return "", err
	}
	v := fmt.Sprintf("map<%s, %s> %s = %d;", keyTypeString, c.ValueTyping, c.Name, c.Tag)
	if c.Comment != "" {
		v = fmt.Sprintf("%s   // %s", v, c.Comment)
	}
	return v, nil
}

func (e Enum) Write(level int) (string, error) {
	v := fmt.Sprintf("%senum %s {\n", indentLevel(level), e.Name)
	if e.AllowAlias {
		v = fmt.Sprintf("%s%soption allow_alias = true;\n", v, indentLevel(level+1))
	}
	for _, enumValue := range e.Values {
		v = fmt.Sprintf("%s%s%s = %d;", v, indentLevel(level+1), enumValue.Name, enumValue.Tag)
		if enumValue.Comment != "" {
			v = fmt.Sprintf("%s   // %s", v, enumValue.Comment)
		}
		v = fmt.Sprintf("%s\n", v)
	}
	v = fmt.Sprintf("%s%s}", v, indentLevel(level))
	return v, nil
}

// VALIDATORS

func (m *Message) Validate() error {
	if m.Name == "" {
		return errors.New("Message name cannot be empty")
	}
	return nil
}

func (s *ScalarField) Validate() error {
	if s.Name == "" {
		return errors.New("Scalar field must have a non-empty name")
	}
	return nil
}

func (m *MapField) Validate() error {
	if m.Name == "" {
		return errors.New("Map field must have a non-empty name")
	}
	if m.KeyTyping < 1 || m.KeyTyping == DOUBLE_TYPE || m.KeyTyping == FLOAT_TYPE || m.KeyTyping == BYTES_TYPE {
		return fmt.Errorf("Map field %s must use a scalar integral or string type for the map key", m.Name)
	}
	if m.ValueTyping < 1 {
		return fmt.Errorf("Map field %s must have a type specified for the map value", m.Name)
	}
	return nil
}

// FORMATTING

func indentLevel(level int) string {
	var buffer bytes.Buffer
	for i := 0; i < level; i++ {
		buffer.WriteString("  ")
	}
	return buffer.String()
}

func toProtobufType(t FieldType) (string, error) {
	switch t {
	case DOUBLE_TYPE:
		return "double", nil
	case FLOAT_TYPE:
		return "float", nil
	case INT32_TYPE:
		return "int32", nil
	case INT64_TYPE:
		return "int64", nil
	case UINT32_TYPE:
		return "uint32", nil
	case UINT64_TYPE:
		return "uint64", nil
	case SINT32_TYPE:
		return "sint32", nil
	case SINT64_TYPE:
		return "sint64", nil
	case FIXED32_TYPE:
		return "fixed32", nil
	case FIXED64_TYPE:
		return "fixed64", nil
	case SFIXED32_TYPE:
		return "sfixed32", nil
	case SFIXED64_TYPE:
		return "sfixed64", nil
	case BOOL_TYPE:
		return "bool", nil
	case STRING_TYPE:
		return "string", nil
	case BYTES_TYPE:
		return "bytes", nil
	default:
		return "", fmt.Errorf("Unrecognized protobuf field type: %d", t)
	}
}
