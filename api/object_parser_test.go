package api

import (
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_ParseObjectId(t *testing.T) {
	var oid *ObjectId
	oidStr := "ff6ccb68859fd52216ec8dadf98d2a00859f5369"
	t1 := objectParserForString(oidStr)
	oid = t1.ParseObjectId()
	util.Assert(t, oid.String() == oidStr)
}
