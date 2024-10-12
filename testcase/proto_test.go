package testcase

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"github.com/ucanme/proto_parser/parser"
	"reflect"
	"testing"
)

func TestMarshal2Json(t *testing.T) {
	convey.Convey("TestMarshal2Json", t, func() {
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
}`

		p, err := parser.Parse([]byte(msgStr))
		fmt.Println("---p----", p)
		convey.So(err, convey.ShouldBeNil)
		m, err := p.Unmarshal2Json("StateMsg", []byte{0x08, 0x03, 0x10, 0x02, 0x18, 0x03})
		convey.So(err, convey.ShouldBeNil)
		data, _ := json.Marshal(m)

		ok, err := DeepEqualJson(string(data), `{"live_state":0,"st_state_detailed":3,"state":2,"state_detailed":3}`)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ok, convey.ShouldBeTrue)
	})
}

func DeepEqualJson(jsonStr1, jsonStr2 string) (ok bool, err error) {
	var (
		json1 interface{}
		json2 interface{}
	)
	if err = json.Unmarshal([]byte(jsonStr1), &json1); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(jsonStr2), &json2); err != nil {
		return
	}
	ok = reflect.DeepEqual(json1, json2)
	return
}
