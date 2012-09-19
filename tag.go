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
    tag string
    tagger *PersonTimestamp
    message string
    typeName string
}

func (t *Tag) String() string {
    const FMT = "object %s\ntype %s\ntag %s\ntagger %s\nmessage %s"
    //TODO
    return fmt.Sprintf(FMT, t.object, t.typeName, t.tag, *(t.tagger), t.message)
}

func (t *Tag) Type() ObjectType {
    return ObjectTag
}

func (t *Tag) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, t.String())
}

func toTag(repo Repository, obj *RawObject) (*Tag, error) {
    var p []byte
    var t *Tag
    var err error
    if p, err = obj.Payload(); err != nil {
        return nil, err
    }
    if t, err = parseTag(p); err != nil {
        fmt.Println("could not parse: ", err)
        return nil, err
    }
    return t, nil // TODO
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
        tag.typeName = line
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
