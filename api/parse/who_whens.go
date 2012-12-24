//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
who_whens.go implements user credentials and timestamps.
*/
package parse

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/token"
	"github.com/jbrukh/ggit/util"
	"strings"
)

var signs []string = []string{
	token.PLUS,
	token.MINUS,
}

// ================================================================= //
// PARSING
// ================================================================= //

func (p *objectParser) parseWhoWhen(marker string) *objects.WhoWhen {
	p.ConsumeString(marker)
	p.ConsumeByte(token.SP)
	user := strings.Trim(p.ReadString(token.LT), string(token.SP))
	email := p.ReadString(token.GT)
	p.ConsumeByte(token.SP)
	seconds := p.ParseInt(token.SP, 10, 64)

	// time zone
	signStr := p.ConsumeStrings(signs)
	var sign int64
	if signStr == token.PLUS {
		sign = 1
	} else if signStr == token.MINUS {
		sign = -1
	} else {
		util.PanicErrf("expecting: +/- sign")
	}

	tzHours := p.ParseIntN(2, 10, 64)
	tzMins := p.ParseIntN(2, 10, 64)
	if tzMins < 0 || tzMins > 59 {
		util.PanicErrf("expecting 00 to 59 for tz minutes")
	}

	// time zone offset in signed minutes
	tz := int(sign * (tzHours*int64(60) + tzMins))

	ww := objects.NewWhoWhen(user, email, seconds, tz)

	return ww
}
