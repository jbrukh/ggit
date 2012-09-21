package ggit

//commits and tags have author-with-timestamp and commiter-with-timestamp fields.
type WhoWhen struct {
    name  string
    email string
    date  string // TODO: turn into date object
}

func (at *WhoWhen) Name() string {
    return at.name
}

func (at *WhoWhen) Email() string {
    return at.email
}

// TODO: turn into date
func (at *WhoWhen) Date() string {
    return at.date
}
