package api

import (
	"testing"
	"fmt"
)

const testTagText = `object a53e5437def11c1eed3a4be1d45fba42e7582b03
type commit
tag tattoo
tagger Michael Bosworth <michael.a.bosworth@gmail.com> 1348058205 -0400

this is not a tag. it is a flying spaghetti monster.`

var testTag string

func init() {
	testTag = fmt.Sprintf("tag %d\000%s", len(testTagText), testTagText)
}

func Test_parseTag(t *testing.T) {
	r := readerForString(testTag)
	p := newObjectParser(r)

	parsed, err := p.ParsePayload()
	if err != nil {
		println(err)
	}

	assertf(t, err == nil, "failed due to error")
	assert(t, parsed != nil, "failed because payload was nil")

	tag, ok := parsed.(*Tag)

	assert(t, ok)
	assert(t, tag.tag == "tattoo")
	assert(t, tag.object.String() == "a53e5437def11c1eed3a4be1d45fba42e7582b03" )
	assert(t, tag.message == "this is not a tag. it is a flying spaghetti monster.")
	assert(t, tag.tagger.Email() == "michael.a.bosworth@gmail.com")
	assert(t, tag.tagger.Name() == "Michael Bosworth")
	assert(t, tag.tagger.Seconds() == 1348058205)
	assert(t, tag.tagger.Offset() == int(-240))
}
