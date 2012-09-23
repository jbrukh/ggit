package ggit

import (
	"testing"
)

var testTreeSha *ObjectId
var testParentSha *ObjectId

func init() {
	testTreeSha, _ = NewObjectIdFromString("e98b3d7be9979411127f93a1b9027c1eb5fe83b4")
	testParentSha, _ = NewObjectIdFromString("8e5c7a9c2f37f315375d26ae8148690f920d2b62")
}

const testData = `tree e98b3d7be9979411127f93a1b9027c1eb5fe83b4
parent 8e5c7a9c2f37f315375d26ae8148690f920d2b62
author Jake Brukhman <brukhman@gmail.com> 1348333582 -0400
committer Jake Brukhman <brukhman@gmail.com> 1348333582 -0400

Structure for WhoWhen.`
const testCommit1 = "commit " + string(len(testData)) + string(NUL) + testData

func Test_parseCommit(t *testing.T) {
	c1 := readerForString(testCommit1)
	p := newObjectParser(c1)

	parsed, err := p.ParsePayload()
	assertf(t, err == nil, "failed due to error")

	c, _ := parsed.(Commit)

	assert(t, c.tree.String() == testTreeSha.String())
	assert(t, c.parents != nil && len(c.parents) != 0)
	assert(t, c.parents[0].String() == testParentSha.String())
	assert(t, c.author.Name() == "Jake Brukhman")
	assert(t, c.author.Email() == "brukhman@gmail.com")
	assert(t, c.author.Seconds() == 1348333582)
	assert(t, c.author.Offset() == -240)
	assert(t, c.committer.Name() == "Jake Brukhman")
	assert(t, c.committer.Email() == "brukhman@gmail.com")
	assert(t, c.committer.Seconds() == 1348333582)
	assert(t, c.committer.Offset() == -240)
	assert(t, c.message == "Structure for WhoWhen.")
}
