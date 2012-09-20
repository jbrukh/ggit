package ggit

import "io"

type Object interface {
    Type() ObjectType

    // write the string representation of
    // this object to the writer
    WriteTo(w io.Writer) (n int, err error)
}
