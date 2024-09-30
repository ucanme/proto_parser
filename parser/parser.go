package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type location struct {
	column int
	line   int
}

type parser struct {
	br             *bufio.Reader
	loc            *location
	eofReached     bool
	prefix         string
	lastColumnRead int
	Messages       map[string]MessageElement
	Enums          map[string]EnumElement
}

var eof = rune(0)

func (p *parser) unread() {
	if p.loc.column == 0 {
		p.loc.line--
		p.loc.column = p.lastColumnRead
	}
	_ = p.br.UnreadRune()
}

func (p *parser) read() rune {
	c, _, err := p.br.ReadRune()
	if err != nil {
		return eof
	}

	p.lastColumnRead = p.loc.column

	if c == '\n' {
		p.loc.line++
		p.loc.column = 0
	} else {
		p.loc.column++
	}
	return c
}

func (p *parser) parse() error {
	p.Enums = map[string]EnumElement{}
	p.Messages = map[string]MessageElement{}
	for {
		p.skipWhitespace()
		if p.eofReached {
			break
		}
		err := p.ReadDelaration(parseCtx{ctxType: fileCtx})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) ReadDelaration(ctx parseCtx) error {
	var err error
	label := p.readWord()
	if label == "message" {
		err = p.readMessage(ctx)
	} else if label == "enum" {
		err = p.readEnum(ctx)
	} else if ctx.ctxType == msgCtx {
		if !ctx.permitsField() {
			return p.errline("fields must be nested")
		}
		err = p.readField(label, "", ctx)
	} else if ctx.ctxType == enumCtx {
		err = p.readOption(ctx)
	} else {
		return errors.New("expected message declaration")
	}
	return err
}

func (p *parser) readMessage(ctx parseCtx) error {
	p.skipWhitespace()
	name, _, err := p.readName()
	if err != nil {
		return err
	}

	me := MessageElement{Name: name, QualifiedName: p.prefix + name, Fields: map[int]FieldElement{}, Messages: map[string]MessageElement{}}
	var previousPrefix = p.prefix
	p.prefix = p.prefix + name + "."
	defer func() {
		p.prefix = previousPrefix
	}()
	p.skipWhitespace()
	if c := p.read(); c != '{' {
		return p.throw('{', c)
	}

	innerCtx := parseCtx{ctxType: msgCtx, obj: &me}
	if err = p.readDeclarationsInLoop(innerCtx); err != nil {
		return err
	}

	if ctx.ctxType == msgCtx {
		parent := ctx.obj.(*MessageElement)
		parent.Messages[me.Name] = me
	} else {
		p.Messages[me.Name] = me
	}
	return nil
}

func (p *parser) readWord() string {
	return p.readWordAdvanced(nil)
}

type enclosure int

const (
	parenthesis enclosure = iota
	bracket
	unenclosed
)

func (p *parser) readName() (string, enclosure, error) {
	var name string
	enc := unenclosed
	c := p.read()
	if c == '(' {
		enc = parenthesis
		name = p.readWord()
		if p.read() != ')' {
			return "", enc, p.errline("Expected ')'")
		}
		p.unread()
	} else if c == '[' {
		enc = bracket
		name = p.readWord()
		if p.read() != ']' {
			return "", enc, p.errline("Expected ']'")
		}
		p.unread()
	} else {
		p.unread()
		name = p.readWord()
	}
	return name, enc, nil
}

func isValidCharInWord(c rune, f func(r rune) bool) bool {
	if isLetter(c) || isDigit(c) || c == '_' || c == '-' || c == '.' {
		return true
	} else if f != nil {
		return f(c)
	}
	return false
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c rune) bool {
	return (c >= '0' && c <= '9')
}

func (p *parser) readWordAdvanced(f func(r rune) bool) string {
	var buf bytes.Buffer
	for {
		c := p.read()
		if isValidCharInWord(c, f) {
			_, _ = buf.WriteRune(c)
		} else {
			p.unread()
			break
		}
	}
	return buf.String()
}

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func (p *parser) skipWhitespace() {
	for {
		c := p.read()
		if c == eof {
			p.eofReached = true
			break
		} else if !isWhitespace(c) {
			p.unread()
			break
		}
	}
}

func (p *parser) errline(msg string, a ...interface{}) error {
	s := fmt.Sprintf(msg, a...)
	return fmt.Errorf(s+" on line: %v", p.loc.line)
}

func (p *parser) throw(expected rune, actual rune) error {
	return p.errcol("Expected %v, but found: %v", strconv.QuoteRune(expected), strconv.QuoteRune(actual))
}

func (p *parser) errcol(msg string, a ...interface{}) error {
	s := fmt.Sprintf(msg, a...)
	return fmt.Errorf(s+" on line: %v, column: %v", p.loc.line, p.loc.column)
}

func (p *parser) readDeclarationsInLoop(ctx parseCtx) error {
	for {
		p.skipWhitespace()
		if p.eofReached {
			return fmt.Errorf("Reached end of input in %v definition (missing '}')", ctx)
		}
		if c := p.read(); c == '}' {
			break
		}
		p.unread()

		if err := p.ReadDelaration(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) readDataType() (DataType, error) {
	name := p.readWord()
	p.skipWhitespace()
	return p.readDataTypeInternal(name)
}

func (p *parser) readDataTypeInternal(name string) (DataType, error) {
	// is it a map type?
	//if name == "map" {
	//	if c := p.read(); c != '<' {
	//		return nil, p.throw('<', c)
	//	}
	//	var err error
	//	var keyType, valueType DataType
	//	keyType, err = p.readDataType()
	//	if err != nil {
	//		return nil, err
	//	}
	//	if c := p.read(); c != ',' {
	//		return nil, p.throw(',', c)
	//	}
	//	p.skipWhitespace()
	//	valueType, err = p.readDataType()
	//	if err != nil {
	//		return nil, err
	//	}
	//	if c := p.read(); c != '>' {
	//		return nil, p.throw('>', c)
	//	}
	//	return MapDataType{keyType: keyType, valueType: valueType}, nil
	//}

	// is it a scalar type?
	sdt, err := NewScalarDataType(name)
	if err == nil {
		return sdt, nil
	}

	if _, ok := p.Enums[name]; ok {
		return EnumDataType{name: name}, nil
	}
	// must be a named type
	return NamedDataType{name: name}, nil
}

func (p *parser) readEnum(ctx parseCtx) error {
	p.skipWhitespace()
	name, _, err := p.readName()
	if err != nil {
		return err
	}
	p.skipWhitespace()
	if c := p.read(); c != '{' {
		return p.throw('{', c)
	}

	ee := EnumElement{Name: name, QualifiedName: p.prefix + name}
	innerCtx := parseCtx{ctxType: enumCtx, obj: &ee}
	if err = p.readDeclarationsInLoop(innerCtx); err != nil {
		return err
	}

	// add enum to the proper parent...
	if ctx.ctxType == msgCtx {
		me := ctx.obj.(*MessageElement)
		me.Enums[ee.Name] = ee
	} else {
		p.Enums[ee.Name] = ee
	}
	return nil
}
