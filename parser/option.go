package parser

import (
	"io"
	"strings"
)

func (p *parser) readOption(ctx parseCtx) error {
	var err error
	var enc enclosure
	oe := OptionElement{}

	p.skipWhitespace()
	if oe.Name, enc, err = p.readName(); err != nil {
		return err
	}
	oe.IsParenthesized = (enc == parenthesis)

	p.skipWhitespace()
	if c := p.read(); c != '=' {
		return p.throw('=', c)
	}
	p.skipWhitespace()

	if p.read() == '"' {
		oe.Value = p.readUntil('"')
	} else {
		p.unread()
		oe.Value = p.readWord()
	}

	p.skipWhitespace()
	if c := p.read(); c != ';' {
		return p.throw(';', c)
	}

	// add option to the proper parent...
	if ctx.ctxType == msgCtx {
		me := ctx.obj.(*MessageElement)
		me.Options = append(me.Options, oe)
	} else if ctx.ctxType == enumCtx {
		ee := ctx.obj.(*EnumElement)
		ee.Options = append(ee.Options, oe)
	}
	return nil
}

func (p *parser) readUntil(delimiter byte) string {
	s, err := p.br.ReadString(delimiter)
	if err == io.EOF {
		p.eofReached = true
	}
	return strings.TrimSuffix(s, string(delimiter))
}
