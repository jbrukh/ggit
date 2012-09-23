package ggit

import (
	"fmt"
	"strings"
)

// Who represents a user that has been using
// (g)git
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

// WhoWhen
type WhoWhen struct {
	Who
	When
}

func (ww *WhoWhen) String() string {
	const format = "%s <%s>"                          // TODO
	return fmt.Sprintf(format, ww.Name(), ww.Email()) // TODO
}

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
