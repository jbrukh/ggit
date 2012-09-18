package ggit

//commits and tags have author-with-timestamp fields.
type AuthorTimestamp struct {
    name  string
    email string
    date  string // TODO: turn into date object
}

func (at *AuthorTimestamp) Name() string {
    return at.name
}

func (at *AuthorTimestamp) Email() string {
    return at.email
}

// TODO: turn into date
func (at *AuthorTimestamp) Date() string {
    return at.date
}
