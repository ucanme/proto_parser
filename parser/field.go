package parser

func (p *parser) readField(label string, documentation string, ctx parseCtx) error {
	fe := FieldElement{Documentation: documentation}
	var err error
	dataTypeStr := label
	if label == required || label == optional || label == repeated {
		fe.Label = label
		p.skipWhitespace()
		dataTypeStr = p.readWord()
	}
	p.skipWhitespace()

	//// figure out the dataType
	if fe.Type, err = p.readDataTypeInternal(dataTypeStr); err != nil {
		return err
	}

	//// figure out the name
	p.skipWhitespace()
	if fe.Name, _, err = p.readName(); err != nil {
		return err
	}
	//
	//// check for equals sign...
	p.skipWhitespace()
	var c rune
	if c = p.read(); c != '=' {
		return p.throw('=', c)
	}

	//
	//// extract the field tag...
	p.skipWhitespace()
	if fe.Tag, err = p.readInt(); err != nil {
		return err
	}
	//
	//// If semicolon is next; we are done. If '[' is next, we must parse options for the field
	if c = p.read(); c != ';' {
		return p.throw(';', c)
	}
	//
	//// add field to the proper parent	...
	if ctx.ctxType == msgCtx {
		me := ctx.obj.(*MessageElement)
		me.Fields[fe.Tag] = fe
	}
	//} else if ctx.ctxType == extendCtx {
	//	ee := ctx.obj.(*ExtendElement)
	//	ee.Fields = append(ee.Fields, fe)
	//} else if ctx.ctxType == oneOfCtx {
	//	oe := ctx.obj.(*OneOfElement)
	//	oe.Fields = append(oe.Fields, fe)
	//}
	return nil
}

func (p *parser) readHeader(data []byte, offset int) (tag int, wireType, size int) {
	header := data[offset]
	wireType = int(0x07 & header)
	tag = int(header >> 3)
	offset++
	size = 1
	return
}

const (
	optional = "optional"
	required = "required"
	repeated = "repeated"
)
