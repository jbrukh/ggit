package api

import (
	"fmt"
	"io"
)

// Formatter is the central hub for formatting ggit objects.
type Format struct {
	Writer io.Writer
}

func (f *Format) Lf() (int, error) {
	return fmt.Fprint(f.Writer, "\n")
}

func (f *Format) Printf(format string, items ...interface{}) (int, error) {
	return fmt.Fprintf(f.Writer, format, items...)
}
