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
package objects

import "time"

// ================================================================= //
// WHO
// ================================================================= //

type Who struct {
	name  string
	email string
}

// seconds: unix time
// offset: time zone offset in signed minutes
func NewWhoWhen(name, email string, seconds int64, offset int) *WhoWhen {
	return &WhoWhen{
		Who{name, email},
		When{seconds, offset},
	}
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
