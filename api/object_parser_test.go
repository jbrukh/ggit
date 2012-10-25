//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_ParseObjectId(t *testing.T) {
	var oid *ObjectId
	oidStr := "ff6ccb68859fd52216ec8dadf98d2a00859f5369"
	t1 := objectParserForString(oidStr)
	oid = t1.ParseOid()
	util.Assert(t, oid.String() == oidStr)
}
