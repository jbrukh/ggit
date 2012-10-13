package api

import (
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_ParseObjectId(t *testing.T) {
	var oid *ObjectId
	t1 := objectParserForString(testOidCrazy)
	oid = t1.ParseObjectId()
	util.Assert(t, oid.String() == testOidCrazy)
}
