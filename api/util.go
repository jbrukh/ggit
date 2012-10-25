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
// HELPERS
// ================================================================= //

// trimLast throws away the last character of a byte slice
func trimLastByte(b []byte) []byte {
	if b == nil || len(b) == 0 {
		return b
	}
	return b[:len(b)-1]
}

func trimLastStr(b []byte) string {
	return string(trimLastByte(b))
}

func trimLast(str string) string {
	if str == "" {
		return str
	}
	return str[:len(str)-1]
}

func trimPrefix(str, prefix string) string {
	for _, v := range prefix {
		if len(str) > 0 && uint8(v) == str[0] {
			str = str[1:]
		} else {
			panic("prefix doesn't match")
		}
	}
	return str
}

// ================================================================= //
// METHODS FOR TESTING
// ================================================================= //

// TODO: move to util

func objectParserForString(str string) *objectParser {
	p := new(objectParser)
	p.buf = readerForString(str)
	return p
}

func parserForBytes(b []byte) *dataParser {
	return &dataParser{
		buf: bufio.NewReader(bytes.NewBuffer(b)),
	}
}

func parserForString(str string) *dataParser {
	return parserForBytes([]byte(str))
}

func readerForString(str string) *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(str))
}
