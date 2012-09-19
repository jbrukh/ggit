package ggit

//commits and tags have author-with-timestamp and commiter-with-timestamp fields.
type PersonTimestamp struct {
    name  string
    email string
    date  string // TODO: turn into date object
}

func (at *PersonTimestamp) Name() string {
    return at.name
}

func (at *PersonTimestamp) Email() string {
    return at.email
}

// TODO: turn into date
func (at *PersonTimestamp) Date() string {
    return at.date
}
