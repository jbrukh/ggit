package ggit

import (
    "fmt"
    "io"
	"bufio"
	"bytes"
)

const (
	tokenObject 	= "object"
	tokenType       = "type"
	tokenTag        = "tag"
	tokenTagger     = "tagger"
)

type Tag struct {
    repo *Repository
    //commit-SHA1
    object *ObjectId
    //tag name
    tag string
    //author with timestamp
    tagger *PersonTimestamp
    //message
    message string
}

func (t *Tag) String() string {
    //TODO
    return ""
}

func (t *Tag) Type() ObjectType {
    return OBJECT_TAG
}

func (t *Tag) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, t.String())
}

func toTag(repo Repository, obj *RawObject) (t *Tag, err error) {
    p, e := obj.Payload()
    if e != nil {
        return nil, e
    }
    // TODO implement the parsing
    fmt.Println(string(p))
    return new(Tag), nil // TODO
}

func parseTag(b []byte) (*Tag, error) {
	p := &dataParser{bufio.NewReader(bytes.NewBuffer(b))}
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
		line := p.ReadString(LF)                  // gets rid of the LF!
		//TODO: what do? do tags ever refer to e.g. other tags, or trees, or is type always "commit"?

		// read the tag
		p.ConsumeString(tokenTag)
		p.ConsumeByte(SP)
		line = p.ReadString(LF)                      // gets rid of the LF!
		tag.tag = line

		// read the tagger
		p.ConsumeString(tokenTagger)
		p.ConsumeByte(SP)
		line = p.ReadString(LF)                      // gets rid of the LF!
		tag.tagger = &PersonTimestamp{line, "", ""} // TODO

		// read the commit message
		p.ConsumeByte(LF)
		tag.message = p.String()
	})
	return tag, err
}
