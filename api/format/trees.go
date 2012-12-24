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
	"github.com/jbrukh/ggit/api/token"
)

// ================================================================= //
// FORMATTING
// ================================================================= //

// Tree formats this tree object into an API-friendly string that is
// the same as the output of git-cat-file tree <tree>.
func (f *formatter) Tree(t *objects.Tree) (int, error) {
	N := 0
	for _, e := range t.Entries() {
		n, err := fmt.Fprintf(f.Writer, "%o %s%s%s", e.Mode(), e.Name(), string(token.NUL), string(e.ObjectId().Bytes()))
		N += n
		if err != nil {
			return N, err
		}
	}
	return N, nil
}

// TreePretty formats this tree object into a human-friendly table
// that is the same as the output of git-cat-file -p <tree>.
func (f *formatter) TreePretty(t *objects.Tree) (int, error) {
	N := 0
	for _, e := range t.Entries() {
		n, err := fmt.Fprintf(f.Writer, "%.6o %s %s\t%s\n", e.Mode(), e.ObjectType(), e.ObjectId(), e.Name())
		N += n
		if err != nil {
			return N, err
		}
	}
	return N, nil
}
