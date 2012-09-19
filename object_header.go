package ggit

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type ObjectHeader struct {
	Type ObjectType
	Size int
}
