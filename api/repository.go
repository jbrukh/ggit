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
	"compress/zlib"
	"fmt"
	"github.com/jbrukh/ggit/util"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	DefaultGitDir     = ".git"
	DefaultObjectsDir = "objects"
	IndexFile         = "index"
	PackedRefsFile    = "packed-refs"
)

// Repository. Currently, this interface is tracking
// the interface of DiskRepository (for the most part).
// However, in the scheme of things, a Repository
// should be a more general interface.
type Repository interface {

	// Destroy will mercilessly and irreparably delete
	// the existing repository.
	Destroy() error

	// TODO: this needs to be replaced with
	// higher level index operations
	Index() (*Index, error)

	// TODO: while this is ok for now, this debug
	// method should not be part of the backend interface
	ObjectIds() ([]*ObjectId, error)

	// Refs returns a list of all refs in the repository.
	// TODO: perhaps replace with a visitor of refs?
	Refs() ([]Ref, error)

	// Ref convert a string ref into a Ref object. The
	// returned object may be a symbolic or concrete ref.
	Ref(spec string) (Ref, error)

	// ObjectFromOid is the fundamental object retrieval
	// operation of a repository. It is the basis for
	// working with any object.
	ObjectFromOid(oid *ObjectId) (Object, error)

	// ObjectFromShortOid provides support for shortened
	// hashes. This functionality is usually tied to the
	// particular kind of backend the repository is using.
	ObjectFromShortOid(short string) (Object, error)
}

// a representation of a git repository
type DiskRepository struct {
	path string
	pr   []Ref
}

// Open a reprository that is located at the given path. The
// path of a repository is always its .git directory. However,
// if the enclosing directory is given, then ggit will
// append the .git directory to the specified path.
func Open(pth string) *DiskRepository {
	p := util.InferGitDir(pth)
	return &DiskRepository{
		path: p,
	}
}

// Destroy is a highly destructive operation that 
// irrevocably destroys the git repository and its
// enclosing directory.
func (repo *DiskRepository) Destroy() error {
	dir, _ := filepath.Split(repo.path)
	return os.RemoveAll(dir)
}

func (repo *DiskRepository) ObjectFromOid(oid *ObjectId) (obj Object, err error) {
	var (
		f  *os.File
		e  error
		rz io.ReadCloser
	)
	if f, e = repo.objectFile(oid); e != nil {
		return nil, e
	}
	defer f.Close() // just in case
	if rz, e = zlib.NewReader(f); e != nil {
		return nil, e
	}
	defer rz.Close()
	file := bufio.NewReader(rz)
	p := newObjectParser(file, oid)
	return p.ParsePayload()
}

func (repo *DiskRepository) ObjectFromShortOid(short string) (Object, error) {
	l := len(short)
	if l < 4 || l > OidHexSize {
		return nil, fmt.Errorf("fatal: Not a valid object name %s", short)
	}

	// don't bother with directories if we know the full SHA
	if l == OidHexSize {
		oid, err := OidFromString(short)
		if err != nil {
			return nil, fmt.Errorf("fatal: Not a valid object name %s", short)
		}
		return repo.ObjectFromOid(oid)
	}

	head, tail := short[:2], short[2:]
	root := path.Join(DefaultGitDir, DefaultObjectsDir, head)
	var matching []*ObjectId
	e := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// root doesn't exist, or there was a problem reading it
		if err != nil {
			return err
		}
		if !info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, tail) {
				if oid, err := OidFromString(head + name); err == nil {
					matching = append(matching, oid)
				}
			}
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	if len(matching) != 1 {
		return nil, fmt.Errorf("fatal: Ambiguous object name %s", short)
	}
	return repo.ObjectFromOid(matching[0])
}

func (repo *DiskRepository) Ref(spec string) (Ref, error) {
	file, e := repo.relativeFile(spec)
	if e == nil {
		defer file.Close()
		p := newRefParser(bufio.NewReader(file), spec)
		return p.parseRef()
	}

	if os.IsNotExist(e) {
		// we can check packed refs now
		// TODO: we can optimize this by caching it
		refs, err := repo.PackedRefs()
		if err != nil {
			return nil, noSuchRefErrf(spec)
		}
		for _, r := range refs {
			if r.Name() == spec {
				return r, nil
			}
		}
		return nil, noSuchRefErrf(spec)
	}

	return nil, e
}

// find all objects and print their ids
func (repo *DiskRepository) ObjectIds() (oids []*ObjectId, err error) {
	objectsRoot := path.Join(repo.path, DefaultObjectsDir)
	oids = make([]*ObjectId, 0)
	//look in each objectsDir and make ObjectIds out of the files there.
	err = filepath.Walk(objectsRoot, func(path string, info os.FileInfo, errr error) error {
		if name := info.Name(); name == "info" || name == "pack" {
			return filepath.SkipDir
		} else if !info.IsDir() {
			hash := filepath.Base(filepath.Dir(path)) + name
			var oid *ObjectId
			if oid, err = OidFromString(hash); err != nil {
				return err
			}
			oids = append(oids, oid)
		}
		return nil
	})
	return
}

func (repo *DiskRepository) Index() (idx *Index, err error) {
	file, e := repo.relativeFile(IndexFile)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	return toIndex(bufio.NewReader(file))
}

func (repo *DiskRepository) PackedRefs() (pr []Ref, err error) {
	if repo.pr == nil {
		file, e := repo.relativeFile(PackedRefsFile)
		if e != nil {
			return nil, e
		}
		defer file.Close()
		p := newRefParser(bufio.NewReader(file), "")
		if pr, e = p.ParsePackedRefs(); e != nil {
			return nil, e
		}
		repo.pr = pr
	}
	return pr, nil
}

func (repo *DiskRepository) LooseRefs() ([]Ref, error) {
	// TODO: figure out a way to decouple this logic
	repoPath := repo.path + "/"
	dir := path.Join(repoPath, "refs")
	refs := make([]Ref, 0)
	err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				spec := trimPrefix(path, repoPath)
				r, e := OidRefFromRef(repo, spec)
				if e != nil {
					return e
				}
				refs = append(refs, &ref{name: spec, oid: r.ObjectId()})
			}
			return nil
		},
	)
	return refs, err
}

func (repo *DiskRepository) Refs() ([]Ref, error) {

	// First, get all the packed refs.
	pr, err := repo.PackedRefs()
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
	dir := path.Join(repo.path, "refs")
	err = filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			// refs are files, so...
			if !f.IsDir() {
				spec := trimPrefix(path, repo.path+"/")
				r, e := OidRefFromRef(repo, spec)
				if e != nil {
					return e
				}
				refs[spec] = &ref{name: spec, oid: r.ObjectId()}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	// collect the refs into a list
	refList := make([]Ref, 0, len(refs))
	for _, v := range refs {
		refList = append(refList, v)
	}
	sort.Sort(refByName(refList))
	return refList, nil
}

// ================================================================= //
// PRIVATE METHODS
// ================================================================= //

// turn an oid into a path relative to the
// git directory of a repository
func (repo *DiskRepository) objectFile(oid *ObjectId) (file *os.File, err error) {
	hex := oid.String()
	path := path.Join(repo.path, DefaultObjectsDir, hex[0:2], hex[2:])
	return os.Open(path)
}

func (repo *DiskRepository) relativeFile(relPath string) (file *os.File, err error) {
	path := path.Join(repo.path, relPath)
	return os.Open(path)
}
