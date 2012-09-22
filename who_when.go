package ggit

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
