package api

import (
	"fmt"
	"github.com/jbrukh/ggit/util"
	"testing"
)

const testTagText = `object a53e5437def11c1eed3a4be1d45fba42e7582b03
type commit
tag tattoo
tagger Michael Bosworth <michael.a.bosworth@gmail.com> 1348058205 -0400

this is not a tag. it is a flying spaghetti monster.`

var testTag string

func withHeader(tagText string) string {
	return fmt.Sprintf("tag %d\000%s", len(tagText), tagText)
}

func init() {
	testTag = withHeader(testTagText)
}

func Test_parseTag(t *testing.T) {
	tag, ok := parseTag(t, testTag)
	util.Assert(t, ok)
	util.Assert(t, tag.name == "tattoo")
	util.Assert(t, tag.object.String() == "a53e5437def11c1eed3a4be1d45fba42e7582b03")
	util.Assert(t, tag.message == "this is not a tag. it is a flying spaghetti monster.")
	util.Assert(t, tag.tagger.Email() == "michael.a.bosworth@gmail.com")
	util.Assert(t, tag.tagger.Name() == "Michael Bosworth")
	util.Assert(t, tag.tagger.Seconds() == 1348058205)
	util.Assert(t, tag.tagger.Offset() == int(-240))
}

func Test_tagString(t *testing.T) {
	tag, ok := parseTag(t, testTag)
	util.Assert(t, ok)
	f := NewStrFormat()
	f.Tag(tag)
	s := f.String()
	tagData := withHeader(s)
	var thereAndBackAgain *Tag
	thereAndBackAgain, ok = parseTag(t, tagData)
	util.Assert(t, ok)
	util.Assert(t, tag.name == thereAndBackAgain.name)
	util.Assert(t, tag.object.String() == thereAndBackAgain.object.String())
	util.Assert(t, tag.message == thereAndBackAgain.message)
	util.Assert(t, tag.tagger.Email() == thereAndBackAgain.tagger.Email())
	util.Assert(t, tag.tagger.Name() == thereAndBackAgain.tagger.Name())
	util.Assert(t, tag.tagger.Seconds() == thereAndBackAgain.tagger.Seconds())
	util.Assert(t, tag.tagger.Offset() == thereAndBackAgain.tagger.Offset())
}

func parseTag(t *testing.T, s string) (tag *Tag, ok bool) {
	r := readerForString(s)
	oid := OidNow("ff6ccb68859fd52216ec8dadf98d2a00859f5369")
	p := newObjectParser(r, oid)

	parsed, err := p.ParsePayload()
	if err != nil {
		println(err)
	}

	util.Assertf(t, err == nil, "failed due to error")
	util.Assert(t, parsed != nil, "failed because payload was nil")

	tag, ok = parsed.(*Tag)
	return
}
