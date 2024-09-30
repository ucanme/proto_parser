package parser

import (
	"bytes"
	"strconv"
)

type FieldElement struct {
	Name          string
	Documentation string
	Options       []OptionElement
	Label         string /* optional, required, repeated, oneof */
	Type          DataType
	Tag           int
}

// OneOfElement is a datastructure which models
// oneof construct share memory, and at most one field can be
// set at any time.
type OneOfElement struct {
	Name          string
	Documentation string
	Options       []OptionElement
	Fields        []FieldElement
}

type ExtensionsElement struct {
	Documentation string
	Start         int
	End           int
}

type ExtendElement struct {
	Name          string
	QualifiedName string
	Documentation string
	Fields        []FieldElement
}

type ReservedRangeElement struct {
	Documentation string
	Start         int
	End           int
}

type MessageElement struct {
	Name               string
	QualifiedName      string
	Documentation      string
	Options            []OptionElement
	Fields             map[int]FieldElement
	Enums              map[string]EnumElement
	Messages           map[string]MessageElement
	OneOfs             []OneOfElement
	ExtendDeclarations []ExtendElement
	Extensions         []ExtensionsElement
	ReservedRanges     []ReservedRangeElement
	ReservedNames      []string
}

type OptionElement struct {
	Name            string
	Value           string
	IsParenthesized bool
}

// EnumConstantElement is a datastructure which models
// the fields within an enum construct. Enum constants can
// also have inline options specified.
type EnumConstantElement struct {
	Name          string
	Documentation string
	Options       []OptionElement
	Tag           int
}

// EnumElement is a datastructure which models
// the enum construct in a protobuf file. Enums are
// defined standalone or as nested entities within messages.
type EnumElement struct {
	Name          string
	QualifiedName string
	Documentation string
	Options       []OptionElement
	EnumConstants []EnumConstantElement
}

func (p *parser) readInt() (int, error) {
	var buf bytes.Buffer
	for {
		c := p.read()
		if isDigit(c) {
			_, _ = buf.WriteRune(c)
		} else {
			p.unread()
			break
		}
	}
	str := buf.String()
	intVal, err := strconv.Atoi(str)
	return intVal, err
}
