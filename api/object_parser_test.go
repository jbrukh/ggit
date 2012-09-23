package api

import (
	"testing"
)

func Test_ParseObjectId(t *testing.T) {
	var oid *ObjectId
	t1 := objectParserForString(testOidCrazy)
	oid = t1.ParseObjectId()
	assert(t, oid.String() == testOidCrazy)
}
