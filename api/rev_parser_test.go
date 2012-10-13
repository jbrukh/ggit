package api

import (
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_parseNumber(t *testing.T) {
	test := func(str string, exp int) {
		p := &revParser{
			spec: str,
		}
		i, err := p.parseNumber()
		util.Assert(t, err == nil && i == exp)
	}

	test("~~", 1)
	test("~1", 1)
	test("~1001~1", 1001)
	test("^^", 1)
	test("^1", 1)
	test("^1001^1", 1001)

}
