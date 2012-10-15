//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"bytes"
	"fmt"
	"io"
)

// Formatter is the central hub for formatting ggit objects.
type Format struct {
	Writer io.Writer
}

type strFormat struct {
	Format
}

func (f *strFormat) String() string {
	return f.Writer.(*bytes.Buffer).String()
}

func (f *strFormat) Reset() {
	f.Writer.(*bytes.Buffer).Reset()
}

func NewStrFormat() *strFormat {
	b := bytes.NewBufferString("")
	return &strFormat{
		Format{b},
	}
}

func (f *Format) Lf() (int, error) {
	return fmt.Fprint(f.Writer, "\n")
}

func (f *Format) Printf(format string, items ...interface{}) (int, error) {
	return fmt.Fprintf(f.Writer, format, items...)
}
