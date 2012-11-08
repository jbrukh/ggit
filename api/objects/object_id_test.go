//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package objects

import (
	"crypto/sha1"
	"fmt"
	"github.com/jbrukh/ggit/util"
	"io"
	"testing"
)

const (
	testOidZeros = "0000000000000000000000000000000000000000"
	testOidOnes  = "1111111111111111111111111111111111111111"
)

const numTests = 100

func getTestOids(size int) []string {
	if size < 2 {
		panic("make size bigger")
	}
	res := make([]string, size)
	// always add the edge cases
	res[0] = testOidOnes
	res[1] = testOidZeros
	for i := 2; i < size; i++ {
		res[i] = util.UniqueHex20()
	}
	return res
}

func Test_OidFromString(t *testing.T) {
	ids := getTestOids(numTests)
	for _, id := range ids {
		oid, err := OidFromString(id)
		util.AssertNoErr(t, err)
		util.AssertEqualString(t, oid.String(), id)
	}
}

func Test_OidFromHash(t *testing.T) {
	h := sha1.New()
	id := OidFromHash(h)
	if id.bytes == nil {
		t.Error("did not initialize bytes properly")
	}

	// now we will test a real hash
	io.WriteString(h, "I have always known that one day I would take this road, but yesterday I did not know it would be today.")
	hashBytes := h.Sum(nil)[0:OidSize]
	id = OidFromHash(h)

	expected := fmt.Sprintf("%x", hashBytes)
	actual := id.String()
	if actual == testOidZeros || expected == testOidZeros || actual != expected {
		t.Error("bad hash initialization: ", expected, " but got ", actual)
	}
}
