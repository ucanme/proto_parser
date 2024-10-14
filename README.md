# proto_parser
A dynamic parser for protobuf messages define and a tranformer for binary data in proto protocal to json string in golang language.

# OverView 
proto_parser is really convient to  convert binary data in proto protocal to  json string, and it is also easy to use.

# How to use
1. pick out  your message define from your proto file and parse it to a paser.
```protobuf
    msgStr :=
	`enum State {
    STATE_ONE= 0;
    STATE_TWO   = 1;
    STATE_THREE = 2;
    STATE_FOUR    = 3;
    STATE_FIVE   = 7;
}

enum LiveState {
    LIVE_ST_NONE=0;
    LIVE_ST_MD=1;
    LIVE_ST_AD= 2;
    LIVE_ST_INVALID=7;
}

enum StateDetailed {
    ST_PARKED   = 0;
    ST_DRVRDY   = 1;
    ST_DRIVING  = 2;
    ST_SWUPDATE = 3;
    ST_CHARGING = 4;
    ST_PWRSWAP  = 5;
    ST_INVALID = 15;
}
enum STStateDetailed {
    ST_DETAILED_NONE    = 0;
    ST_DETAILED_MD      = 1;
    ST_DETAILED_NAD     = 2;
    ST_DETAILED_RD      = 3;
    ST_DETAILED_INVALID = 7;
}
message StateMsg {
    StateDetailed state_detailed = 1;
    State state                  = 2;
    STStateDetailed st_state_detailed = 3;
    LiveState live_state = 4;
}
 p, err := parser.Parse([]byte(msgStr))
```
2. unmarshal binary data to json string
```go
    data := []byte{0x08, 0x03, 0x10, 0x02, 0x18, 0x03} // some proto binary data
    jsonStr, err := p.Unmarshal2Json(msgName, data)
    fmt.Println(string(jsonStr))
```

3. unmarshal result
```json
    {
        "state_detailed": 3,
        "state": 2,
        "st_state_detailed": 0,
        "live_state": 3
    }
```
