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
    // ReadRawObject reads a raw object from the backend
    ReadRawObject(oid *ObjectId) (obj *RawObject, err error)

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

func (r *DiskRepository) ReadRawObject(oid *ObjectId) (obj *RawObject, err error) {
    file, err := r.objectFile(oid)
    if err != nil {
        return
    }
    defer file.Close()

    var zr io.ReadCloser
    if zr, err = zlib.NewReader(file); err == nil {
        defer zr.Close()
        obj = new(RawObject)
        _, err = io.Copy(obj, zr)
    }
    return
}

func (r *DiskRepository) ReadObject(oid *ObjectId) (obj Object, err error) {
    rawObj, err := r.ReadRawObject(oid)
    if err != nil {
        return
    }

    h, err := rawObj.Header()
    if err != nil {
        return
    }

    switch h.Type {
    case OBJECT_BLOB:
        return toBlob(r, rawObj)
    case OBJECT_TREE:
        return toTree(r, rawObj)
    case OBJECT_COMMIT:
        return toCommit(r, rawObj)
    default:
        panic("unsupported type")
    }
    return
}

func (r *DiskRepository) Index() (idx *Index, err error) {
    file, e := r.indexFile()
    if e != nil {
        return
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
