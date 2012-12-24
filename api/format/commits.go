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

// TODO: the return values of this method are broken
func (f *formatter) Commit(c *objects.Commit) (int, error) {
	// tree
	fmt.Fprintf(f.Writer, "tree %s\n", c.Tree())

	// parents
	for _, p := range c.Parents() {
		fmt.Fprintf(f.Writer, "parent %s\n", p)
	}

	// author
	sf := NewStrFormat()
	sf.WhoWhen(c.Author())
	fmt.Fprintf(f.Writer, "author %s\n", sf.String())
	sf.Reset()
	sf.WhoWhen(c.Committer())
	fmt.Fprintf(f.Writer, "committer %s\n", sf.String())

	// commit message
	fmt.Fprintf(f.Writer, "\n%s", c.Message())
	return 0, nil // TODO TODO
}
