package api

// import (
// 	"testing"
// )

const testTag = `object a53e5437def11c1eed3a4be1d45fba42e7582b03
type commit
tag tattoo
tagger Michael Bosworth <michael.a.bosworth@gmail.com> 1348058205 -0400

this is not a tag. it is a tattoo on your ass with your own name on it. get it? go fuck yourself.`

// func Test_parseTag(t *testing.T) {
// 	r := readerForString(testTag)

// 	tag, _ := parseTag(nil, &objectHeader{ObjectTag, len(testTag)}, r)

// 	assert(t, tag.tag == "tattoo")
// 	assert(t, tag.object.String() == "a53e5437def11c1eed3a4be1d45fba42e7582b03" )
// 	assert(t, tag.message == "this is not a tag. it is a tattoo on your ass with your own name on it. get it? go fuck yourself.")
// 	assert(t, tag.tagger.Email() == "michael.a.bosworth@gmail.com")
// 	assert(t, tag.tagger.Name() == "Michael Bosworth")
// 	assert(t, tag.tagger.Seconds() == 1348058205)
// 	assert(t, tag.tagger.Offset() == int(-240))
// }
