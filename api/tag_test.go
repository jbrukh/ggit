package api

import (
	"fmt"
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
	assert(t, ok)
	assert(t, tag.tag == "tattoo")
	assert(t, tag.object.String() == "a53e5437def11c1eed3a4be1d45fba42e7582b03")
	assert(t, tag.message == "this is not a tag. it is a flying spaghetti monster.")
	assert(t, tag.tagger.Email() == "michael.a.bosworth@gmail.com")
	assert(t, tag.tagger.Name() == "Michael Bosworth")
	assert(t, tag.tagger.Seconds() == 1348058205)
	assert(t, tag.tagger.Offset() == int(-240))
}

func Test_tagString(t *testing.T) {
	tag, ok := parseTag(t, testTag)
	assert(t, ok)
	f := NewStrFormat()
	f.Tag(tag)
	s := f.String()
	tagData := withHeader(s)
	var thereAndBackAgain *Tag
	thereAndBackAgain, ok = parseTag(t, tagData)
	assert(t, ok)
	assert(t, tag.tag == thereAndBackAgain.tag)
	assert(t, tag.object.String() == thereAndBackAgain.object.String())
	assert(t, tag.message == thereAndBackAgain.message)
	assert(t, tag.tagger.Email() == thereAndBackAgain.tagger.Email())
	assert(t, tag.tagger.Name() == thereAndBackAgain.tagger.Name())
	assert(t, tag.tagger.Seconds() == thereAndBackAgain.tagger.Seconds())
	assert(t, tag.tagger.Offset() == thereAndBackAgain.tagger.Offset())
}

func parseTag(t *testing.T, s string) (tag *Tag, ok bool) {
	r := readerForString(s)
	p := newObjectParser(r)

	parsed, err := p.ParsePayload()
	if err != nil {
		println(err)
	}

	assertf(t, err == nil, "failed due to error")
	assert(t, parsed != nil, "failed because payload was nil")

	tag, ok = parsed.(*Tag)
	return
}
