//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
package format

import "github.com/jbrukh/ggit/api/objects"

// ================================================================= //
// FORMATTING
// ================================================================= //

// The "plumbing" string output of an Object. This output can be used
// to reproduce the contents or SHA1 hash of an Object.
func (f *Format) Object(o objects.Object) (int, error) {
	/* TODO: consider that this isn't extensible.
	 * Outside of the format package, new Object types cannot add their
	 * own formatting here. */
	switch t := o.(type) {
	case *objects.Blob:
		return f.Blob(t)
	case *objects.Tree:
		return f.Tree(t)
	case *objects.Commit:
		return f.Commit(t)
	case *objects.Tag:
		return f.Tag(t)
	}
	panic("unknown object")
}

// The pretty string output of an Object. This format is not necessarily
// of use as an api call; it is for humans.
func (f *Format) ObjectPretty(o objects.Object) (int, error) {
	switch t := o.(type) {
	case *objects.Blob:
		return f.Blob(t)
	case *objects.Tree:
		return f.TreePretty(t)
	case *objects.Commit:
		return f.Commit(t)
	case *objects.Tag:
		return f.TagPretty(t)
	}
	panic("unknown object")
}
