//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
)

const (
	markerObject = "object"
	markerType   = "type"
	markerTag    = "tag"
	markerTagger = "tagger"
)

// ================================================================= //
// PARSING
// ================================================================= //

func (p *objectParser) parseTag() *objects.Tag {
	p.ResetCount()

	// read the object id
	p.ConsumeString(markerObject)
	p.ConsumeByte(SP)
	target := p.ParseOid()
	p.ConsumeByte(LF)

	// read object type
	p.ConsumeString(markerType)
	p.ConsumeByte(SP)
	t := objects.ObjectType(p.ConsumeStrings(objectTypes))
	p.ConsumeByte(LF)

	// read the tag name
	p.ConsumeString(markerTag)
	p.ConsumeByte(SP)
	name := p.ReadString(LF) // gets rid of the LF!

	// read the tagger
	tagger := p.parseWhoWhen(markerTagger)
	p.ConsumeByte(LF)

	// read the commit message
	p.ConsumeByte(LF)
	msg := p.String()

	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}

	return objects.NewTag(p.oid, target, t, p.hdr, name, msg, tagger)
}

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) Tag(t *objects.Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.Object())
	fmt.Fprintf(f.Writer, "type %s\n", t.ObjectType())
	fmt.Fprintf(f.Writer, "tag %s\n", t.Name())
	sf := NewStrFormat()
	sf.WhoWhen(t.Tagger())
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.Message())
	return 0, nil // TODO
}

func (f *Format) TagPretty(t *objects.Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.Object())
	fmt.Fprintf(f.Writer, "type %s\n", t.ObjectType())
	fmt.Fprintf(f.Writer, "tag %s\n", t.Name())
	sf := NewStrFormat()
	sf.WhoWhenDate(t.Tagger()) // git-cat-file -p displays full dates for tags
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.Message())
	return 0, nil // TODO
}
