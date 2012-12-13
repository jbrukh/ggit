//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
package format

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
)

// ================================================================= //
// OBJECT FORMATTER
// ================================================================= //

// Blob formats the contents of the blog as a string
// for output to the screen.
func (f *Format) Blob(b *objects.Blob) (int, error) {
	return fmt.Fprintf(f.Writer, "%s", string(b.Data()))
}
