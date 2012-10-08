package api

import (
	"testing"
)

func Test_parseNumber(t *testing.T) {
	p := &revParser{
		spec: "~1",
	}

	i, err := p.parseNumber()
	assert(t, err == nil && i == 1)
}
