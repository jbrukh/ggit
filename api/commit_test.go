package api

import (
	"fmt"
	"github.com/jbrukh/ggit/util"
	"testing"
)

var testTreeSha *ObjectId
var testParentSha *ObjectId

func init() {
	testTreeSha, _ = OidFromString("e98b3d7be9979411127f93a1b9027c1eb5fe83b4")
	testParentSha, _ = OidFromString("8e5c7a9c2f37f315375d26ae8148690f920d2b62")
}

const testData = `tree e98b3d7be9979411127f93a1b9027c1eb5fe83b4
parent 8e5c7a9c2f37f315375d26ae8148690f920d2b62
author Jake Brukhman <brukhman@gmail.com> 1348333582 -0400
committer Jake Brukhman <brukhman@gmail.com> 1348333582 -0400

Structure for WhoWhen.`

var testCommit1 string

func init() {
	testCommit1 = fmt.Sprintf("commit %d\000%s", len(testData), testData)
}

func Test_parseCommit(t *testing.T) {
	c1 := readerForString(testCommit1)
	oid, _ := OidFromString(testOidCrazy)
	p := newObjectParser(c1, oid)

	parsed, err := p.ParsePayload()
	if err != nil {
		fmt.Println("error was: ", err.Error())
	}
	util.Assertf(t, err == nil, "failed due to error")
	util.Assert(t, parsed != nil, "parsed is nul")

	c, ok := parsed.(*Commit)
	util.Assert(t, ok)
	util.Assert(t, c.tree.String() == testTreeSha.String())
	util.Assert(t, c.parents != nil && len(c.parents) != 0)
	util.Assert(t, c.parents[0].String() == testParentSha.String())
	util.Assert(t, c.author.Name() == "Jake Brukhman")
	util.Assert(t, c.author.Email() == "brukhman@gmail.com")
	util.Assert(t, c.author.Seconds() == 1348333582)
	util.Assert(t, c.author.Offset() == -240)
	util.Assert(t, c.committer.Name() == "Jake Brukhman")
	util.Assert(t, c.committer.Email() == "brukhman@gmail.com")
	util.Assert(t, c.committer.Seconds() == 1348333582)
	util.Assert(t, c.committer.Offset() == -240)
	util.Assert(t, c.message == "Structure for WhoWhen.")

}
