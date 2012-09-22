package ggit

import (
	"bufio"
	"compress/zlib"
	"errors"
	"io"
	"os"
	"path"
)

const DEFAULT_GIT_DIR = ".git"
const INDEX_FILE = "index"

// A Backend supports storage of arbitrary Git
// objects without particular regard of the technical
// specifics. The backend can deliver a RawObject
// by id (it is a read-only key-value store.)
type Backend interface {
	// Read an arbitrary object from the backend
	ReadObject(oid *ObjectId) (o Object, err error)
}

type Repository interface {
	Backend
	Index() (idx *Index, err error)
}

// a representation of a git repository
type DiskRepository struct {
	path string
}

// open a reprository that is located at the given path
func Open(path string) (r *DiskRepository, err error) {
	// check that repo is valid
	if !validateRepo(path) {
		return nil, errors.New("not a valid repo")
	}
	r = &DiskRepository{
		path: path,
	}
	return
}

// closing operations for a repository
func (r *DiskRepository) Close() {
}

func (r *DiskRepository) ReadObject(oid *ObjectId) (obj Object, err error) {
	var (
		f  *os.File
		e  error
		rz io.ReadCloser
	)
	if f, e = r.objectFile(oid); e != nil {
		return nil, e
	}
	defer f.Close() // just in case

	if rz, e = zlib.NewReader(f); e != nil {
		return nil, e
	}
	defer rz.Close()

	file := bufio.NewReader(rz)

	h, err := parseObjectHeader(file)
	if err != nil {
		return
	}

	switch h.Type {
	case ObjectBlob:
		return parseBlob(r, h, file)
	case ObjectTree:
		return parseTree(r, h, file)
	case ObjectCommit:
		return parseCommit(r, h, file)
	case ObjectTag:
		return parseTag(r, h, file)
	default:
		panic("unsupported type")
	}
	return
}

func (r *DiskRepository) Index() (idx *Index, err error) {
	file, e := r.indexFile()
	if e != nil {
		return nil, e
	}
	defer file.Close()
	return toIndex(bufio.NewReader(file))
}

// IndexFile returns an open git index file. It is up to the
// caller to close this resource.
func (r *DiskRepository) indexFile() (file *os.File, err error) {
	path := path.Join(r.path, INDEX_FILE)
	return os.Open(path)
}

// turn an oid into a path relative to the
// git directory of a repository
func (r *DiskRepository) objectFile(oid *ObjectId) (file *os.File, err error) {
	hex := oid.String()
	path := path.Join(r.path, "objects", hex[0:2], hex[2:])
	return os.Open(path)
}

// validate a repository path to make sure it has
// the right format and that it exists
func validateRepo(path string) bool {
	// TODO
	return true
}
