package api

import (
	"fmt"
	"strings"
	"time"
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

func (w *When) Date() string {
	t := time.Unix(w.seconds, int64(0))
	// standard time: Mon Jan 2 15:04:05 -0700 MST 2006
	return t.Format("Mon Jan 2 15:04:05 2006")

}

func (f *Format) WhoWhen(ww *WhoWhen) (int, error) {
	sign := ""
	if ww.offset < 0 {
		sign = MINUS
	}
	offset := abs(ww.Offset())
	hours := int(offset / 60)
	minutes := offset - hours*60
	return fmt.Fprintf(f.Writer, "%s <%s> %s %s%02d%02d", ww.Name(), ww.Email(), ww.Date(), sign, hours, minutes)
}

// WhoWhen
type WhoWhen struct {
	Who
	When
}

func (ww *WhoWhen) String() string {
	// TODO: move this to Format
	const format = "%s <%s> %d %s"
	offset := ww.Offset()
	hours := int(offset / 60)
	minutes := fmt.Sprintf("%d", (offset - (hours * 60)))
	if len(minutes) == 1 {
		//pad with 0
		minutes = "0" + minutes
	}
	zone := fmt.Sprintf("%d%s", hours, minutes)
	if len(zone) == 4 {
		//pad hour with 0
		zone = string(zone[0]) + "0" + string(zone[1:])
	}
	return fmt.Sprintf(format, ww.Name(), ww.Email(), ww.Seconds(), zone)
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
