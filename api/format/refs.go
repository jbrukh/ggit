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
// REF FORMATTING
// ================================================================= //

func (f *Format) Ref(r objects.Ref) (int, error) {
	_, rf := r.Target() // symbolic or oid
	return fmt.Fprintf(f.Writer, "%s %s", rf, r.Name())
}

// TODO: come up with a better name for this
func (f *Format) Deref(r objects.Ref) (int, error) {
	return fmt.Fprintf(f.Writer, "%s %s^{}", r.Commit(), r.Name())
}
