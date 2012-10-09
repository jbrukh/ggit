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

type Tag struct {
	//commit-SHA1
	object     *ObjectId
	tag        string
	tagger     *WhoWhen
	message    string
	size       int
	objectType ObjectType
	oid        *ObjectId
}

func (t *Tag) Type() ObjectType {
	return ObjectTag
}

func (t *Tag) Size() int {
	return t.size
}

func (t *Tag) ObjectId() *ObjectId {
	return t.oid
}

func (t *Tag) Object() *ObjectId {
	return t.object
}

// ================================================================= //
// OBJECT PARSER
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
	tag.tag = p.ReadString(LF) // gets rid of the LF!

	// read the tagger
	tag.tagger = p.parseWhoWhen(markerTagger)
	p.ConsumeByte(LF)

	// read the commit message
	p.ConsumeByte(LF)
	tag.message = p.String()
	tag.size = p.hdr.Size

	if p.Count() != p.hdr.Size {
		panicErr("payload doesn't match prescibed size")
	}

	return tag
}

// ================================================================= //
// OBJECT FORMATTER
// ================================================================= //

func (f *Format) Tag(t *Tag) (int, error) {
	return fmt.Fprintf(f.Writer, "object %s\ntype %s\ntag %s\ntagger %s\n\n%s",
		t.object, t.objectType, t.tag, t.tagger.String(), t.message)
}
