package ggit

import (
	"path"
	"errors"
    "os"
	"io"
	"compress/zlib"
)

// a representation of a git repository
type Repository struct {
	path string
    wdir string
}

type ObjectDatabase interface {
    ReadRawObject(oid *ObjectId) (o *RawObject, err error)
    ReadBlob(oid *ObjectId) (b *Blob, err error)
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

func (r *Repository) ReadRawObject(oid *ObjectId) (o *RawObject, err error) {
    var file *os.File
    path := path.Join(r.path, objectPath(oid))

    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    var zr io.ReadCloser
    if zr, err = zlib.NewReader(file); err != nil {
        return
    }
    defer zr.Close()

    o = new(RawObject)
    _, err = io.Copy(o, zr)
    return
}

func (r *Repository) ReadBlob(oid *ObjectId) (b *Blob, err error) {
    var o *RawObject
    if o, err = r.ReadRawObject(oid); err != nil {
        return
    }
    return &Blob{*o}, nil
}


// turn an oid into a path relative to the
// git directory of a repository
func objectPath(oid *ObjectId) string {
	hex := oid.String()
	return path.Join("objects", hex[0:2], hex[2:])
}

// validate a repository path to make sure it has
// the right format and that it exists
func validateRepo(path string) bool {
    return true
}
