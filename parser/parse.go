package parser

import (
	"bufio"
	"bytes"
)

func Parse(data []byte) (p *parser, err error) {
	r := bytes.NewReader(data)
	br := bufio.NewReader(r)
	loc := location{line: 1, column: 0}
	p = &parser{
		br:  br,
		loc: &loc,
	}
	p.parse()
	return p, nil
}
