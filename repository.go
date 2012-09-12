package ggit

import (
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
    ReadRawObject(oid *ObjectId) (o *RawObject, err error)
}

type ObjectReader interface {
    // ReadBlob returns a Blob object representing the
    // blob in question. If there is no such blob (or the
    // id does not refer to a blob), an error is returned.
    ReadBlob(oid *ObjectId) (b *Blob, err error)
    ReadTree(oid *ObjectId) (t *Tree, err error)
    ReadCommit(oid *ObjectId) (c *Commit, err error)
}

// a representation of a git repository
type Repository struct {
    path string
}

// open a reprository that is located at the given path
func Open(path string) (r *Repository, err error) {
    // check that repo is valid
    if !validateRepo(path) {
        return nil, errors.New("not a valid repo")
    }
    r = &Repository{
        path: path,
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
    if zr, err = zlib.NewReader(file); err == nil {
        defer zr.Close()
        o = new(RawObject)
        _, err = io.Copy(o, zr)
    }
    return
}

// ReadRawObjectHeader reads the header information of an object in the repository
func (r *Repository) ReadRawObjectHeader(oid *ObjectId) (h *ObjectHeader, err error) {
    // TODO: could do special parsing to get header info only
    if obj, err := r.ReadRawObject(oid); err == nil {
        h, err = obj.Header()
    }
    return
}

func (r *Repository) ReadObject(oid *ObjectId) (obj Object, err error) {
    rawObj, err := r.ReadRawObject(oid)
    if err != nil {
        return
    }

    h, err := rawObj.Header()
    if err != nil {
        return
    }

    // TODO: fix the double lookup
    switch h.Type {
    case OBJECT_BLOB:
        return r.ReadBlob(oid)
    case OBJECT_TREE:
        return r.ReadTree(oid)
    case OBJECT_COMMIT:
        return r.ReadCommit(oid)
    default:
        panic("unsupported type")
    }
    return
}

// ReadBlob obtains a Blob object from
// the Backend
func (r *Repository) ReadBlob(oid *ObjectId) (b *Blob, err error) {
    rawObj, err := r.ReadRawObject(oid)
    if err != nil {
        return
    }

    b = &Blob{
        RawObject: *rawObj,
        repo:      r,
    }
    return
}

// ReadTree obtains a Tree object from
// the Backend
func (r *Repository) ReadTree(oid *ObjectId) (t *Tree, err error) {
    rawObj, err := r.ReadRawObject(oid)
    if err != nil {
        return
    }
    rawTree := newRawTree(rawObj)
    t, err = rawTree.ParseTree()
    if err != nil {
        return
    }
    // associate
    t.repo = r
    return
}

// ReadCommit obtains a Commit object from
// the Backend
func (r *Repository) ReadCommit(oid *ObjectId) (c *Commit, err error) {
    rawObj, err := r.ReadRawObject(oid)
    if err != nil {
        return
    }
    rc := newRawCommit(rawObj)
    c, err = rc.ParseCommit()
    if err != nil {
        return
    }
    // associate
    c.repo = r
    return
}

// IndexFile returns an open git index file. It is up to the
// caller to close this resource.
func (r *Repository) IndexFile() (file *os.File, err error) {
    path := path.Join(r.path, INDEX_FILE)
    return os.Open(path)
}

// turn an oid into a path relative to the
// git directory of a repository
// TODO: return a file here, just like in IndexFile()
func objectPath(oid *ObjectId) string {
    hex := oid.String()
    return path.Join("objects", hex[0:2], hex[2:])
}

// validate a repository path to make sure it has
// the right format and that it exists
func validateRepo(path string) bool {
    return true
}
