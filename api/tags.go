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
// TAGS
// ================================================================= //

type Tag struct {
	hdr     *objects.ObjectHeader // the size of the tag
	oid     *objects.ObjectId     // the oid of the tag itself
	name    string                // the tag name
	object  *objects.ObjectId     // the object this tag is pointing at
	otype   objects.ObjectType    // the object type
	tagger  *WhoWhen              // the tagger
	message string                // the tag message
}

func (t *Tag) Header() *objects.ObjectHeader {
	return t.hdr
}

func (t *Tag) ObjectId() *objects.ObjectId {
	return t.oid
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) Object() *objects.ObjectId {
	return t.object
}

func (t *Tag) ObjectType() objects.ObjectType {
	return t.otype
}

func (t *Tag) Tagger() *WhoWhen {
	return t.tagger
}

func (t *Tag) Message() string {
	return t.message
}

// ================================================================= //
// PARSING
// ================================================================= //

func (p *objectParser) parseTag() *Tag {
	tag := new(Tag)
	tag.oid = p.oid
	p.ResetCount()

	// read the object id
	p.ConsumeString(markerObject)
	p.ConsumeByte(SP)
	tag.object = p.ParseOid()
	p.ConsumeByte(LF)

	// read object type
	p.ConsumeString(markerType)
	p.ConsumeByte(SP)
	tag.otype = objects.ObjectType(p.ConsumeStrings(objectTypes))
	p.ConsumeByte(LF)

	// read the tag name
	p.ConsumeString(markerTag)
	p.ConsumeByte(SP)
	tag.name = p.ReadString(LF) // gets rid of the LF!

	// read the tagger
	tag.tagger = p.parseWhoWhen(markerTagger)
	p.ConsumeByte(LF)

	// read the commit message
	p.ConsumeByte(LF)
	tag.message = p.String()
	tag.hdr = p.hdr

	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}

	return tag
}

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) Tag(t *Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.object)
	fmt.Fprintf(f.Writer, "type %s\n", t.otype)
	fmt.Fprintf(f.Writer, "tag %s\n", t.name)
	sf := NewStrFormat()
	sf.WhoWhen(t.tagger)
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.message)
	return 0, nil // TODO
}

func (f *Format) TagPretty(t *Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.object)
	fmt.Fprintf(f.Writer, "type %s\n", t.otype)
	fmt.Fprintf(f.Writer, "tag %s\n", t.name)
	sf := NewStrFormat()
	sf.WhoWhenDate(t.tagger) // git-cat-file -p displays full dates for tags
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.message)
	return 0, nil // TODO
}
