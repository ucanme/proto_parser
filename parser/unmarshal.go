package parser

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

var offset int

func (p *parser) Unmarshal(msgName string, data []byte) {
	for _, msg := range p.Messages {
		if msg.Name == msgName {
			for _, field := range msg.Fields {
				p.unmarshal(field, msg, data)
			}
		}
	}
}

func (p *parser) unmarshal(field FieldElement, msg MessageElement, data []byte) {
	switch field.Type.Category() {
	case ScalarDataTypeCategory:
		if field.Label == "repeated" {
			p.ReadMsgHeader(data)
		}
		fieldName := field.Name
		p.Read(data)
		fmt.Println("field name: ", fieldName)
	case NamedDataTypeCategory:
		p.ReadMsgHeader(data)
		for _, message := range msg.Messages {
			if message.Name == field.Type.Name() {
				p.Read(data)
				for _, subField := range message.Fields {
					p.unmarshal(subField, message, data)
				}
			}
		}
		break
	}
}

func (p *parser) ReadMsgHeader(data []byte) {
	header := data[offset]
	wireType := 0x07 & header
	tag := header >> 3
	fmt.Println("tag: ", tag, "wireType", wireType)
	offset++
	x, n := proto.DecodeVarint(data[offset:])
	offset = offset + n
	fmt.Println("---msgLen----", x)
}

func (p *parser) Read(data []byte) {
	fmt.Printf("---data---%0b", data[offset])
	header := data[offset]
	wireType := 0x07 & header
	tag := header >> 3
	fmt.Println("tag: ", tag, "wireType", wireType)
	offset++
	switch wireType {
	case 0:
		x, n := proto.DecodeVarint(data[offset:])
		fmt.Println("----basic----", x, n)
		offset = offset + n
	case 1:
	case 2:
		x, n := proto.DecodeVarint(data[offset:])
		offset += n
		//fmt.Println("----len----", x, n)
		st := data[offset : offset+int(x)]
		offset = offset + int(x)
		fmt.Println("----str----", string(st))
		return
	case 3:
	case 4:

	case 5:

	}
	//如果碰到repeated 或者msg或者string，bytes 读取长度字段
}
