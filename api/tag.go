package api

import (
	"fmt"
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
	hdr        ObjectHeader // the size of the tag
	object     *ObjectId    // the object this tag is pointing at
	name       string       // the tag name
	tagger     *WhoWhen     // the tagger
	message    string       // the tag message
	objectType ObjectType   // the object type
	oid        *ObjectId    // the oid of the tag itself
}

func (t *Tag) Header() ObjectHeader {
	return t.hdr
}

func (t *Tag) ObjectId() *ObjectId {
	return t.oid
}

func (t *Tag) Object() *ObjectId {
	return t.object
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
	tag.object = p.ParseObjectId()
	p.ConsumeByte(LF)

	// read object type
	p.ConsumeString(markerType)
	p.ConsumeByte(SP)
	tag.objectType = ObjectType(p.ConsumeStrings(objectTypes))
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
		panicErr("payload doesn't match prescibed size")
	}

	return tag
}

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) Tag(t *Tag) (int, error) {
	fmt.Fprintf(f.Writer, "object %s\n", t.object)
	fmt.Fprintf(f.Writer, "type %s\n", t.objectType)
	fmt.Fprintf(f.Writer, "tag %s\n", t.name)
	sf := NewStrFormat()
	sf.WhoWhenDate(t.tagger) // git-cat-file -p displays full dates for tags
	fmt.Fprintf(f.Writer, "tagger %s\n\n", sf.String())
	fmt.Fprintf(f.Writer, "%s", t.message)
	return 0, nil // TODO
}
