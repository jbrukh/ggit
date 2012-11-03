//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
util.go implements miscellaneous common functions.
*/
package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

const DefaultGitDir = ".git"

func InferGitDir(pth string) string {
	_, file := filepath.Split(pth)
	if file != DefaultGitDir {
		return path.Join(pth, DefaultGitDir)
	}
	return pth
}

// IsValidRepo validates a repository path to make sure it has
// the right format and that it exists.	
func IsValidRepo(pth string) bool {
	p := InferGitDir(pth)
	if _, e := os.Stat(p); e != nil {
		return false
	}
	// TODO: may want to do other checks here...
	return true
}

// UniqueHex16 generates a random 16-character
// hexadecimal string.
func UniqueHex16() string {
	return UniqueHex20()[:16]
}

// UniqueHex20 generates a random 20-character
// hexadecimal string.
func UniqueHex20() string {
	buf := make([]byte, 20)
	io.ReadFull(rand.Reader, buf)
	return fmt.Sprintf("%x", buf)
}

// IsDigit returns true if and only if the parameter
// is a digit from 0 to 9.
func IsDigit(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// trimLast throws away the last character of a byte slice
func TrimLastByte(b []byte) []byte {
	if b == nil || len(b) == 0 {
		return b
	}
	return b[:len(b)-1]
}

func TrimLastStr(b []byte) string {
	return string(TrimLastByte(b))
}

func TrimLast(str string) string {
	if str == "" {
		return str
	}
	return str[:len(str)-1]
}

func TrimPrefix(str, prefix string) string {
	for _, v := range prefix {
		if len(str) > 0 && uint8(v) == str[0] {
			str = str[1:]
		} else {
			panic("prefix doesn't match")
		}
	}
	return str
}
