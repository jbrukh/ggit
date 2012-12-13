//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
format.go implements the ggit Format object, which contains formatting
functions for various objects that we want to print as output. Each formatting
function is implemented in the file relevant to the object that it formats.

Format prints to an io.Writer, or optionally, one can print a string by
using a strFormat object and calling String() on it.
*/
package format

import (
	"bytes"
	"fmt"
	"io"
)

// Format is a collection of formatting
// methods for various ggit objects we 
// wish to format and output.
type Format struct {
	Writer io.Writer
}

// strFormat does everything that Format
// does, except its output destination is
// a string that you can obtain by calling
// the String() method.
type StrFormat struct {
	Format
}

func NewStrFormat() *StrFormat {
	b := bytes.NewBufferString("")
	return &StrFormat{
		Format{b},
	}
}

func (f *StrFormat) String() string {
	return f.Writer.(*bytes.Buffer).String()
}

func (f *StrFormat) Reset() {
	f.Writer.(*bytes.Buffer).Reset()
}

// Lf prints a line feed to the underlying
// io.Writer of the Format.
func (f *Format) Lf() (int, error) {
	return fmt.Fprint(f.Writer, "\n")
}

func (f *Format) Printf(format string, items ...interface{}) (int, error) {
	return fmt.Fprintf(f.Writer, format, items...)
}
