package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	Parse([]byte(`
	  message Result {
    		string url = 1;
    		string title = 2;
    		repeated string snippets = 3;
     }
	message SearchResponse {
     repeated Result results = 1;
    int32 age = 2;
} `))
}
