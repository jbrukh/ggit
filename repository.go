package ggit

import (
    "bytes"
    "compress/zlib"
    "errors"
    "fmt"
    "io"
    "os"
    "path"
)

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
}

// a representation of a git repository
type Repository struct {
    path string
}

// open a repository that is located at the given path
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
    if zr, err = zlib.NewReader(file); err != nil {
        return
    }
    defer zr.Close()

    o = new(RawObject)
    _, err = io.Copy(o, zr)
    return
}

// ReadBlob obtains a Blob object 
func (r *Repository) ReadBlob(oid *ObjectId) (b *Blob, err error) {
    var o *RawObject
    if o, err = r.ReadRawObject(oid); err != nil {
        return
    }
    b = &Blob{o}
    // TODO: check validity!
    return b, nil
}

// TODO: this is currently broken
func (r *Repository) ReadTree(oid *ObjectId) error {
    o, err := r.ReadRawObject(oid)
    if err != nil {
        return err
    }
    payload, err := o.Payload()
    if err != nil {
        return err
    }
    b := bytes.NewBuffer(payload)
    for {
        modeName, err := b.ReadString('\000')
        if err != nil {
            break
        }
        fmt.Printf("%v\n", modeName)
        bts := b.Next(20)
        hsh := NewObjectIdFromBytes(bts)
        fmt.Printf("sha: %s\n", hsh.String())
        if err != nil {
            break
        }
    }
    return nil
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
