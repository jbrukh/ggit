//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"hash"
)

const (
	NUL = '\000'
	SP  = ' '
	LF  = '\n'
	LT  = '<'
	GT  = '>'
	TAB = '\t'
)

const (
	PLUS  = "+"
	MINUS = "-"
)

var signs []string = []string{
	PLUS,
	MINUS,
}

// the hash object used to build
// hashes of our objects
var sha hash.Hash = sha1.New()

// produce the SHA1 hash for any Object.
/*func makeHash(o Object) hash.Hash {
	sha.Reset()
	kind := string(o.Type())
	content := o.String()
	len := len([]byte(content)) + 1
	toHash := []byte(kind + " " + fmt.Sprint(len) + "\000" + content + "\n")
	sha.Write(toHash)
	return sha
}*/

// get the first OID_SZ of the hash
func getHash(h hash.Hash) []byte {
	return h.Sum(nil)[0:OidSize]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ================================================================= //
// METHODS FOR TESTING
// ================================================================= //

// TODO: move to util

func objectParserForString(str string) *objectParser {
	return &objectParser{
		objectIdParser: *newObjectIdParser(readerForString(str)),
	}
}

func readerForString(str string) *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(str))
}
