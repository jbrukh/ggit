package ggit

import (
	"bufio"
	"fmt"
	"io"
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
	const FMT = "object %s\ntype %s\ntag %s\ntagger %s\nmessage %s"
	//TODO
	return fmt.Sprintf(FMT, t.object, t.objectType, t.tag, *(t.tagger), t.message)
}

func (t *Tag) Type() ObjectType {
	return ObjectTag
}

func (t *Tag) Size() int {
	return t.size
}

func (t *Tag) WriteTo(w io.Writer) (n int, err error) {
	return io.WriteString(w, t.String())
}

func parseTag(repo Repository, h *objectHeader, buf *bufio.Reader) (*Tag, error) {
	p := &dataParser{buf}
	tag := new(Tag)
	err := dataParse(func() {
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
		tag.tagger = parseWhoWhen(p, markerTagger)
		p.ConsumeByte(LF)

		// read the commit message
		p.ConsumeByte(LF)
		tag.message = p.String()
	})
	tag.repo = repo
	tag.size = h.Size
	return tag, err
}
