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
// FORMATTING
// ================================================================= //

func (f *Format) ObjectId(oid *objects.ObjectId) (int, error) {
	return fmt.Fprint(f.Writer, oid.String())
}
