package api

import "io"

type Object interface {
	// Type returns the ObjectType of this Object
	Type() ObjectType

	// Size returns the size of the payload TODO
	Size() int

	// write the string representation of
	// this object to the writer
	WriteTo(w io.Writer) (n int, err error)
}
