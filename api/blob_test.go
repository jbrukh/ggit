package api

import (
	//	"fmt"
	"github.com/jbrukh/ggit/test"
	"testing"
)

const varRepo = "../test/var"

func Test_readSimpleBlobs(t *testing.T) {
	const blob2 = "00750edc07d6415dcc07ae0351e9397b0222b7ba"
	dir, err := test.Repo(varRepo, "../test/cases/linear_history.sh")
	assertNoErr(t, err)

	repo := Open(dir)

	var o Object
	o, err = repo.ObjectFromShortOid(blob2)
	assertNoErr(t, err)

	assert(t, o.Header().Type() == ObjectBlob)
	assert(t, o.Header().Size() == 2)
	assert(t, o.(*Blob).ObjectId().String() == blob2)

	err = repo.Destroy()
	assertNoErr(t, err)
}
