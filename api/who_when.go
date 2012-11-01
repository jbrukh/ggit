//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
who_when.go implements user credentials and timestamps.
*/
package api

import (
	"fmt"
	"strings"
	"time"
)

// ================================================================= //
// WHO
// ================================================================= //

type Who struct {
	name  string
	email string
}

func (w *Who) Name() string {
	return w.name
}

func (w *Who) Email() string {
	return w.email
}

// ================================================================= //
// WHEN
// ================================================================= //

type When struct {
	seconds int64 // seconds since epoch
	offset  int   // timezone offset in minutes
}

func (w *When) Seconds() int64 {
	return w.seconds
}

func (w *When) Offset() int {
	return w.offset
}

func (w *When) Date() string {
	t := time.Unix(w.seconds, int64(0))
	// standard time: Mon Jan 2 15:04:05 -0700 MST 2006
	return t.Format("Mon Jan 2 15:04:05 2006")

}

// ================================================================= //
// WHO WHEN
// ================================================================= //

type WhoWhen struct {
	Who
	When
}

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) WhoWhenDate(ww *WhoWhen) (int, error) {
	return fmt.Fprintf(f.Writer, "%s <%s> %s %s", ww.Name(), ww.Email(), ww.Date(), zone(ww.offset))
}

func (f *Format) WhoWhen(ww *WhoWhen) (int, error) {
	return fmt.Fprintf(f.Writer, "%s <%s> %d %s", ww.Name(), ww.Email(), ww.Seconds(), zone(ww.offset))
}

// ================================================================= //
// PARSING
// ================================================================= //

func (p *objectParser) parseWhoWhen(marker string) *WhoWhen {
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
		panicErrf("expecting: +/- sign")
	}

	tzHours := p.ParseIntN(2, 10, 64)
	tzMins := p.ParseIntN(2, 10, 64)
	if tzMins < 0 || tzMins > 59 {
		panicErrf("expecting 00 to 59 for tz minutes")
	}

	// time zone offset in signed minutes
	tz := int(sign * (tzHours*int64(60) + tzMins))

	ww := &WhoWhen{
		Who{user, email},
		When{seconds, tz},
	}
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
