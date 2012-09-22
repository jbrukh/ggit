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
	timestamp int64
	offset    int // timezone offset in minutes
}

func (w *When) Timestamp() int64 {
	return w.timestamp
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

func parseWhoWhen(p *dataParser, marker string) *WhoWhen {
	p.ConsumeString(marker)
	p.ConsumeByte(SP)
	user := strings.Trim(p.ReadString(LT), string(SP))
	email := p.ReadString(GT)
	p.ConsumeByte(SP)
	seconds := p.ParseInt(SP, 10, 64)
	ww := &WhoWhen{
		Who{user, email},
		When{seconds, 0},
	} // TODO
	return ww
}
