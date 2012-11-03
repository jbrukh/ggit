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
	"fmt"
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
func MakeHash(o Object) (hash.Hash, error) {
	sha.Reset()
	kind := string(o.Header().Type())
	f := NewStrFormat()
	if _, err := f.Object(o); err != nil {
		return nil, err
	}
	content := f.String()
	len := len([]byte(content))
	value := kind + " " + fmt.Sprint(len) + "\000" + content
	toHash := []byte(value)
	sha.Write(toHash)
	return sha, nil
}

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
