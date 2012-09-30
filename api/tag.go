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
	repo Repository
	//commit-SHA1
	object     *ObjectId
	tag        string
	tagger     *WhoWhen
	message    string
	size       int
	objectType ObjectType
}

func (t *Tag) String() string {
	const FMT = "object %s\ntype %s\ntag %s\ntagger %s\n\n%s"
	return fmt.Sprintf(FMT, t.object, t.objectType, t.tag, t.tagger.String(), t.message)
}

func (t *Tag) Type() ObjectType {
	return ObjectTag
}

func (t *Tag) Size() int {
	return t.size
}

// ================================================================= //
// OBJECT PARSER
// ================================================================= //

func (p *objectParser) parseTag() *Tag {
	tag := new(Tag)
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

func (f *Formatter) FormatTag(t *Tag) {
	fmt.Fprint(f.W, t.String())
}
