package api

import (
	"bufio"
	"compress/zlib"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
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
	ReadRef(refPath string) (Ref, error)
	ReadRefs() ([]Ref, error)
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
			if name := info.Name(); name == "info" || name == "pack" {
				return filepath.SkipDir
			} else if !info.IsDir() {
				hash := filepath.Base(filepath.Dir(path)) + name
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

// ReadRef turns a path into a ggit Ref object. By path here
// we mean a ggit refspec, a path relative to the git directory.
func (r *DiskRepository) ReadRef(path string) (re Ref, err error) {
	file, e := r.refFile(path)
	if e != nil {
		return nil, e
	}
	defer file.Close()

	// a ref c
	p := newRefParser(bufio.NewReader(file))
	var oid *ObjectId
	err = safeParse(func() {
		// is it a symbolic ref?
		if p.PeekString(4) == "ref:" {
			p.ConsumeString("ref: ")
			symbolic := p.ReadString(LF)
			ref, e := r.ReadRef(symbolic)
			if e != nil {
				panicErr(e.Error())
	}
			re = &NamedRef{ref.ObjectId(), path}
		} else {
			// parse the object id
			oid = p.ParseObjectId()
			p.ConsumeByte(LF)
			re = &NamedRef{oid, path}
		}
	})
	return re, err
}

func (r *DiskRepository) ReadRefs() ([]Ref, error) {

	// First, get all the packed refs.
	pr, err := r.PackedRefs()
	if err != nil {
		return nil, err
	}

	// Refs will be stores in a map by their symbolic name.
	refs := make(map[string]Ref)
	for _, ref := range pr {
		refs[ref.Name()] = ref
	}

	// Now let's walk loose refs and collect them to supercede
	// the packed refs. It is worth it to note here that
	// packed refs may contain outdated references because
	// they are updated lazily.
	dir := path.Join(r.path, "refs")
	err = filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			// refs are files, so...
			if !f.IsDir() {
				refspec := trimPrefix(path, r.path+"/")
				ref, e := r.ReadRef(refspec)
				if e != nil {
					return e
				}
				refs[refspec] = ref
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	// collect the refs into a list
	l := len(refs)
	refList := make([]Ref, 0, l)
	for _, v := range refs {
		refList = append(refList, v)
	}
	sort.Sort(refByName(refList))
	return refList, nil
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
