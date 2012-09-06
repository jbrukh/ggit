package ggit

import (
	"path"
	"os"
	"io"
	"bytes"
	"strings"
	"errors"
	"strconv"
	"compress/zlib"
)

// a representation of a git repository
type Repository struct {
	path string
    wdir string
}

// open a repository that is located at the given path
func Open(path string) (r *Repository, err error) {
	// check that repo is valid
    if !validateRepo(path) {
        return nil, errors.New("not a valid repo")
    }
    r = &Repository{
		path: path,
        wdir: path, // TODO: this can be generalized
	}
	return
}

// closing operations for a repository
func (r *Repository) Close() {
}

/*
func (r *Repository) ReadBlob(oid *ObjectId) (b *Blob, err error) {
	var file *os.File
	filePath := path.Join(r.path, objectPath(oid))

	file, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	
	buf := new(bytes.Buffer)
	zr, err := zlib.NewReader(file)
	defer zr.Close()
	
	if _, err = io.Copy(buf, zr); err != nil {
		return
	}
	header, err := buf.ReadString('\000')
	
	// remove last character
	header = header[:len(header)-1]
	splt := strings.Split(header, " ")
	var l int
	l, err = strconv.Atoi(splt[1])
	if splt[0] != "blob" || err != nil {
		return nil, errors.New("Not a blob")
	}
	
	bts := buf.Bytes()
	if len(bts) != l {
		return nil, errors.New("Size mismatch")
	}
	
	b = &Blob{
		bytes: bts,
	}
	return
}
*/

// turn an oid into a path relative to the
// git directory of a repository
func objectPath(oid *ObjectId) string {
	hex := oid.String()
	return path.Join("objects", hex[0:2], hex[2:])
}

// validate a repository path to make sure it has
// the right format and that it exists
func validateRepo(path string) bool {
}
