//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
format.go implements the ggit formatter object, which contains formatting
functions for various objects that we want to print as output. Each formatting
function is implemented in the file relevant to the object that it formats.

formatter prints to an io.Writer, or optionally, one can print a string by
using a StrFormat object and calling String() on it.
*/
package format

import (
	"bytes"
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"io"
)

// PrettyFormat prints a human-friendly String representation of certain ggit objects.
type PrettyFormat interface {
	Format
	ObjectPretty(objects.Object) (int, error)
	TagPretty(*objects.Tag) (int, error)
	TreePretty(*objects.Tree) (int, error)
}

// Format prints an API-friendly string representation of a ggit object. The output should
// be in a form such that the object it represents can readily be reconstructed from it.
type Format interface {
	Object(objects.Object) (int, error)
	ObjectId(*objects.ObjectId) (int, error)
	Blob(*objects.Blob) (int, error)
	Commit(*objects.Commit) (int, error)
	Tag(*objects.Tag) (int, error)
	Tree(*objects.Tree) (int, error)
	// the ref's target and name
	Ref(objects.Ref) (int, error)
	// the ref's commit and name
	Deref(objects.Ref) (int, error)
	// the WhoWhen's name, email, seconds (unix time), and time zone offset
	WhoWhen(*objects.WhoWhen) (int, error)
	// the WhoWhen's name, email, date, and time zone offset
	WhoWhenDate(*objects.WhoWhen) (int, error)
	Lf() (int, error)
	Printf(format string, items ...interface{}) (int, error)
}

// formatter is a collection of formatting
// methods for various ggit objects we 
// wish to format and output.
type formatter struct {
	Writer io.Writer
}

// strFormat does everything that formatter
// does, except its output destination is
// a string that you can obtain by calling
// the String() method.
type StrFormat struct {
	formatter
}

func NewFormat(writer io.Writer) Format {
	return &formatter{writer}
}

func NewPrettyFormat(writer io.Writer) PrettyFormat {
	return &formatter{writer}
}

func NewStrFormat() *StrFormat {
	b := bytes.NewBufferString("")
	return &StrFormat{
		formatter{b},
	}
}

func (f *StrFormat) String() string {
	return f.Writer.(*bytes.Buffer).String()
}

func (f *StrFormat) Reset() {
	f.Writer.(*bytes.Buffer).Reset()
}

// Lf prints a line feed to the underlying
// io.Writer of the formatter.
func (f *formatter) Lf() (int, error) {
	return fmt.Fprint(f.Writer, "\n")
}

func (f *formatter) Printf(format string, items ...interface{}) (int, error) {
	return fmt.Fprintf(f.Writer, format, items...)
}
