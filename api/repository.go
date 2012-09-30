package api

import (
	"bufio"
	"compress/zlib"
	"errors"

	"io"

	"os"
	"path"
	"path/filepath"
)

const (
	DefaultGitDir     = ".git"
	DefaultObjectsDir = "objects"
	IndexFile         = "index"
	PackedRefsFile    = "packed-refs"
)

// A Backend supports storage of arbitrary Git
// objects without particular regard of the technical
// specifics. The backend can deliver a RawObject
// by id (it is a read-only key-value store.)
type Backend interface {
	// Read an arbitrary object from the backend
	ReadObject(oid *ObjectId) (o Object, err error)
	ObjectIds() (oids []ObjectId, err error)
}

type Repository interface {
	Backend
	Index() (idx *Index, err error)
	PackedRefs() (pr PackedRefs, err error)
	ReadRef(refPath string) (*NamedRef, error)
	ReadRefs() ([]*NamedRef, error)
	PeelRef(symbolic string) (*ObjectId, error)
}

// a representation of a git repository
type DiskRepository struct {
	path string
	pr   PackedRefs
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
	p := newObjectParser(file)

	return p.ParsePayload()
}

//find all objects and print their ids
func (r *DiskRepository) ObjectIds() (oids []ObjectId, err error) {
	objectsRoot := path.Join(r.path, DefaultObjectsDir)
	oids = make([]ObjectId, 0)
	//look in each objectsDir and make ObjectIds out of the files there.
	err = filepath.Walk(objectsRoot, func(path string, info os.FileInfo, errr error) error {
		if info.IsDir() {
			return nil
		}
		switch parent := filepath.Base(filepath.Dir(path)); parent {
		case "info":
			return nil
		case "pack":
			return nil
		default:
			hash := parent + info.Name()
			var oid *ObjectId
			if oid, err = NewObjectIdFromString(hash); err != nil {
				return err
			}
			oids = append(oids, *oid)
		}
		return nil
	})
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

func (r *DiskRepository) PackedRefs() (pr PackedRefs, err error) {
	file, e := r.packedRefsFile()
	if e != nil {
		return nil, e
	}
	defer file.Close()
	p := newRefParser(bufio.NewReader(file))
	if pr, e = p.ParsePackedRefs(); e != nil {
		return nil, e
	}
	return pr, nil
}

func (r *DiskRepository) ReadRef(refPath string) (re *NamedRef, err error) {
	// TODO: validate ref
	file, e := r.refFile(refPath)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	p := newRefParser(bufio.NewReader(file))
	var oid *ObjectId
	if oid, err = p.ParseRefFile(); e != nil {
		return nil, e
	}
	return &NamedRef{
		oid,
		refPath,
	}, nil
}

func (r *DiskRepository) ReadRefs() ([]*NamedRef, error) {
	refs := make([]*NamedRef, 0)
	dir := path.Join(r.path, "refs")
	err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				file, e := os.Open(path)
				if e != nil {
					return e
				}
				p := newRefParser(bufio.NewReader(file))
				var oid *ObjectId
				if oid, err = p.ParseRefFile(); e != nil {
					return e
				}
				re := &NamedRef{
					oid,
					trimPrefix(path, r.path+"/"),
				}
				refs = append(refs, re)
			}
			return nil
		},
	)
	return refs, err
}

func (r *DiskRepository) PeelRef(symbolic string) (*ObjectId, error) {
	return nil, nil // TODO
}

// ================================================================= //
// PRIVATE METHODS
// ================================================================= //

// IndexFile returns an open git index file. It is up to the
// caller to close this resource.
func (r *DiskRepository) indexFile() (file *os.File, err error) {
	path := path.Join(r.path, IndexFile)
	return os.Open(path)
}

// turn an oid into a path relative to the
// git directory of a repository
func (r *DiskRepository) objectFile(oid *ObjectId) (file *os.File, err error) {
	hex := oid.String()
	path := path.Join(r.path, DefaultObjectsDir, hex[0:2], hex[2:])
	return os.Open(path)
}

// packedRefsFile returns an open git packed refs file. It is the
// responsibility of the caller to close it.
func (r *DiskRepository) packedRefsFile() (file *os.File, err error) {
	path := path.Join(r.path, PackedRefsFile)
	return os.Open(path)
}

func (r *DiskRepository) refFile(refPath string) (file *os.File, err error) {
	path := path.Join(r.path, refPath)
	return os.Open(path)
}

// validate a repository path to make sure it has
// the right format and that it exists
func validateRepo(path string) bool {
	// TODO
	return true
}
