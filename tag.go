package ggit

import (
	"bufio"
	"fmt"
	"io"
)

const (
	tokenObject = "object"
	tokenType   = "type"
	tokenTag    = "tag"
	tokenTagger = "tagger"
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
		// read the tree line
		p.ConsumeString(tokenObject)
		p.ConsumeByte(SP)
		tag.object = p.ParseObjectId()
		p.ConsumeByte(LF)

		// read the tagger
		p.ConsumeString(tokenType)
		p.ConsumeByte(SP)
		line := p.ReadString(LF) // gets rid of the LF!
		tag.objectType = ObjectType(line)

		// read the tag
		p.ConsumeString(tokenTag)
		p.ConsumeByte(SP)
		line = p.ReadString(LF) // gets rid of the LF!
		tag.tag = line

		// read the tagger
		p.ConsumeString(tokenTagger)
		p.ConsumeByte(SP)
		line = p.ReadString(LF) // gets rid of the LF!
		tag.tagger = &WhoWhen{  // TODO TODO TODO TODO TODO !!!
			Who{line, ""},
			When{0, 0},
		}

		// read the commit message
		p.ConsumeByte(LF)
		tag.message = p.String()
	})
	tag.repo = repo
	tag.size = h.Size
	return tag, err
}
