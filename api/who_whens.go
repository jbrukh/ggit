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
package api

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
	"strings"
)

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) WhoWhenDate(ww *objects.WhoWhen) (int, error) {
	return fmt.Fprintf(f.Writer, "%s <%s> %s %s", ww.Name(), ww.Email(), ww.Date(), zone(ww.Offset()))
}

func (f *Format) WhoWhen(ww *objects.WhoWhen) (int, error) {
	return fmt.Fprintf(f.Writer, "%s <%s> %d %s", ww.Name(), ww.Email(), ww.Seconds(), zone(ww.Offset()))
}

// ================================================================= //
// PARSING
// ================================================================= //

func (p *objectParser) parseWhoWhen(marker string) *objects.WhoWhen {
	p.ConsumeString(marker)
	p.ConsumeByte(SP)
	user := strings.Trim(p.ReadString(LT), string(SP))
	email := p.ReadString(GT)
	p.ConsumeByte(SP)
	seconds := p.ParseInt(SP, 10, 64)

	// time zone
	signStr := p.ConsumeStrings(signs)
	var sign int64
	if signStr == PLUS {
		sign = 1
	} else if signStr == MINUS {
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

// ================================================================= //
// UTIL
// ================================================================= //

func zone(offset int) string {
	sign := ""
	if offset < 0 {
		sign = MINUS
		offset = -offset
	} else {
		sign = PLUS
	}
	hours := int(offset / 60)
	minutes := offset - hours*60
	return fmt.Sprintf("%s%02d%02d", sign, hours, minutes)
}
