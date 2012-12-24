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

func (f *formatter) Tag(t *objects.Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.Object())
	fmt.Fprintf(f.Writer, "type %s\n", t.ObjectType())
	fmt.Fprintf(f.Writer, "tag %s\n", t.Name())
	sf := NewStrFormat()
	sf.WhoWhen(t.Tagger())
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.Message())
	return 0, nil // TODO
}

func (f *formatter) TagPretty(t *objects.Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.Object())
	fmt.Fprintf(f.Writer, "type %s\n", t.ObjectType())
	fmt.Fprintf(f.Writer, "tag %s\n", t.Name())
	sf := NewStrFormat()
	sf.WhoWhenDate(t.Tagger()) // git-cat-file -p displays full dates for tags
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.Message())
	return 0, nil // TODO
}
