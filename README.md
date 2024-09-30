# proto_parser
A dynamic parser for protobuf messages define and a tranformer for binary data in proto protocal to json string in golang language.

# OverView 
proto_parser is really convient to  convert binary data in proto protocal to  json string, and it is also easy to use.

# How to use
1. pick out  your message define from your proto file and parse it to a paser.
```
    msgStr := ``
	p, err := parser.Parse([]byte(msgStr))
```
2. unmarshal binary data to json string
```
    data := []byte{} // some proto binary data
    jsonStr, err := p.Unmarshal2Json(msgName, data)
```
