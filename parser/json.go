package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func (p *parser) Unmarshal2Json(msgName string, data []byte) (result map[string]interface{}, err error) {
	result = map[string]interface{}{}
	messageElement, exist := p.Messages[msgName]
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse fail reason:%v", r)
		}
	}()
	if !exist {
		err = fmt.Errorf("msgName %s not exist", msgName)
		return
	}
	innerCtx := parseCtx{ctxType: fileCtx, obj: &messageElement}

	result, _, err = p.readMsg(innerCtx, messageElement, data, 0)
	return result, nil
}

func (p *parser) readMsg(ctx parseCtx, messageElement MessageElement, data []byte, offset int) (result map[string]interface{}, offsetInt int, err error) {
	result = make(map[string]interface{})
	var tmp interface{}
	var length uint64
	var size int
	if ctx.ctxType == msgCtx {
		length, size = proto.DecodeVarint(data[offset:])
		offset += size
	}
	endOffset := offset + int(length)
	innerCtx := parseCtx{ctxType: msgCtx, obj: &messageElement}
	var retMap map[string]interface{}
	msgFieldIndex := 0
	for {
		if offset >= len(data)-1 {
			break
		}

		if ctx.ctxType == msgCtx {
			if offset >= endOffset {
				break
			}
		}

		tag, wireType, size := p.readHeader(data, offset)
		offset += size

		field := messageElement.Fields[tag]
		msgFieldIndex++
		label := field.Label
		if label == "repeated" {
			tmp, offset = p.ReadRepeatField(innerCtx, data, wireType, field, offset)
			if field.Type.Category() == NamedDataTypeCategory {
				if _, ok := result[field.Name].([]interface{}); ok {
					result[field.Name] = append(result[field.Name].([]interface{}), tmp)
				} else {
					result[field.Name] = []interface{}{tmp}
				}
			} else {
				result[field.Name] = tmp
			}
		} else {
			retMap, offset = p.ReadField(innerCtx, data, wireType, field, offset)
			for key, val := range retMap {
				result[key] = val
			}
		}
	}

	//填充默认值
	for _, field := range messageElement.Fields {
		if _, ok := result[field.Name]; !ok {
			if _, ok := result[field.Name]; !ok {
				if field.Type.Category() == ScalarDataTypeCategory {
					result[field.Name] = field.Type.DefaultValue()
				} else if field.Type.Category() == NamedDataTypeCategory {
					if _, ok := p.Enums[field.Type.Name()]; ok {
						result[field.Name] = 0
					} else if _, ok := p.Messages[field.Type.Name()]; ok {
						result[field.Name] = map[string]string{}
					}

				} else if field.Type.Category() == MapDataTypeCategory {
					result[field.Name] = field.Type.DefaultValue()
				}
			}
		}
	}

	return result, offset, nil
}

func (p *parser) ReadRepeatField(ctx parseCtx, data []byte, wireType int, element FieldElement, offset int) (result interface{}, offsetInt int) {
	result = make(map[string]interface{})
	var err error

	if element.Type.Category() == NamedDataTypeCategory {
		if _, ok := p.Messages[element.Type.Name()]; ok {
			result, offset, err = p.readMsg(ctx, p.Messages[element.Type.Name()], data, offset)
			if err != nil {
				return
			}
		}
		if _, ok := p.Enums[element.Type.Name()]; ok {
			list := []interface{}{}
			length, size := proto.DecodeVarint(data[offset:])
			offset += size
			endOffset := offset + int(length)
			for offset < endOffset {
				item, size := proto.DecodeVarint(data[offset:])
				offset += size
				list = append(list, item)
			}
			return list, offset
		}

		return result, offset
	}

	if element.Type.Category() == ScalarDataTypeCategory || element.Type.Category() == EnumDataTypeCategory {
		list := []interface{}{}
		length, size := proto.DecodeVarint(data[offset:])
		offset += size
		endOffset := offset + int(length)
		for offset < endOffset {
			if element.Type.Name() == "uint32" || element.Type.Category() == EnumDataTypeCategory {
				item, size := proto.DecodeVarint(data[offset:])
				offset += size
				list = append(list, item)
			}
		}
		return list, offset
	}

	return result, offset
}

func (p *parser) ReadField(ctx parseCtx, data []byte, wireType int, element FieldElement, offset int) (result map[string]interface{}, offsetInt int) {
	result = make(map[string]interface{})
	var retMap map[string]interface{}
	if wireType == 0 {
		item, size := proto.DecodeVarint(data[offset:])
		offset += size
		result[element.Name] = item
	}

	if wireType == 2 {
		if element.Type.Name() == "string" || element.Type.Name() == "bytes" {
			datalen, size := proto.DecodeVarint(data[offset:])
			offset += size
			result[element.Name] = string(data[offset:(offset + int(datalen))])
			offset = offset + int(datalen)
			return result, offset
		}

		obj := ctx.obj.(*MessageElement)
		messageElement := obj.Messages[element.Type.Name()]
		innerCtx := parseCtx{ctxType: msgCtx, obj: &messageElement}
		retMap, offset, _ = p.readMsg(innerCtx, messageElement, data, offset)
		for key, val := range retMap {
			result[key] = val
		}
		for _, field := range messageElement.Fields {
			if _, ok := result[field.Name]; !ok {
				if field.Type.Category() == ScalarDataTypeCategory {
					result[field.Name] = field.Type.DefaultValue()
				} else if field.Type.Category() == NamedDataTypeCategory {
					if _, ok := p.Enums[field.Type.Name()]; ok {
						result[field.Name] = 0
					} else if _, ok := p.Messages[field.Type.Name()]; ok {
						result[field.Name] = map[string]string{}
					}

				} else if field.Type.Category() == MapDataTypeCategory {
					result[field.Name] = field.Type.DefaultValue()
				}
			}
		}
	}

	if wireType == 5 {
		if element.Type.Name() == "float" {
			var f float32
			binary.Read(bytes.NewBuffer(data[offset:offset+4]), binary.LittleEndian, &f)
			offset = offset + 4
			result[element.Name] = f
		}
	}

	return result, offset
}

func (p *parser) IsEnumFieldType(ctx parseCtx, messageElement MessageElement, element FieldElement) bool {
	if element.Type.Category() == EnumDataTypeCategory {
		return true
	}
	if _, exist := messageElement.Enums[element.Type.Name()]; exist {
		return true
	}
	if _, exist := p.Enums[element.Type.Name()]; exist {
		return true
	}
	return false
}
