package api

import (
	"bufio"
	"compress/zlib"
	"errors"
	"fmt"
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

// Repository. Currently, this interface is tracking
// the interface of DiskRepository (for the most part).
// However, in the scheme of things, a Repository
// should be a more general interface.
type Repository interface {
	// TODO: this needs to be replaced with
	// higher level index operations
	Index() (*Index, error)

	// TODO: loose or packed refs may be irrelevant
	// at this level of abstractions, probably should
	// remove from here. For instance, packed refs
	// are meant to compensate for lots of disk reads
	// but such optimization may be irrelevant for
	// repos with distributed cache backends.
	LooseRefs() ([]Ref, error)
	PackedRefs() ([]Ref, error)

	// TODO: while this is ok for now, this debug
	// method should not be part of the backend interface
	ObjectIds() ([]ObjectId, error)

	// Refs returns a list of all refs in the repository.
	// TODO: perhaps replace with a visitor of refs?
	Refs() ([]Ref, error)

	// Ref convert a string ref into a Ref object. The
	// returned object may be a symbolic or concrete ref.
	Ref(spec string) (Ref, error)

	ObjectFromOid(oid *ObjectId) (Object, error)

	RevParse(name string) (Object, error)
}

// a representation of a git repository
type DiskRepository struct {
	path string
	pr   []Ref
}

// open a reprository that is located at the given path
func Open(path string) (*DiskRepository, error) {
	// check that repo is valid
	if !validateRepo(path) {
		return nil, errors.New("not a valid repo")
	}
	return &DiskRepository{
		path: path,
	}, nil
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

func (repo *DiskRepository) Ref(spec string) (Ref, error) {
	file, e := repo.relativeFile(spec)
	if e == nil {
		defer file.Close()
		p := newRefParser(bufio.NewReader(file), spec)
		return p.parseRef()
	}

	if os.IsNotExist(e) {
		// we can check packed refs now
		// TODO: we can optimize this with a trie
		refs, err := repo.PackedRefs()
		if err != nil {
			return nil, err
		}
		refs = FilterRefs(refs, FilterRefPattern(spec))
		if len(refs) != 1 {
			return nil, fmt.Errorf("Ambiguous or missing ref: %s", spec)
		}
		return refs[0], nil
	}
	return nil, e
}

//find all objects and print their ids
func (repo *DiskRepository) ObjectIds() (oids []ObjectId, err error) {
	objectsRoot := path.Join(repo.path, DefaultObjectsDir)
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

func (repo *DiskRepository) Index() (idx *Index, err error) {
	file, e := repo.relativeFile(IndexFile)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	return toIndex(bufio.NewReader(file))
}

func (repo *DiskRepository) PackedRefs() (pr []Ref, err error) {
	file, e := repo.relativeFile(PackedRefsFile)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	p := newRefParser(bufio.NewReader(file), "")
	if pr, e = p.ParsePackedRefs(); e != nil {
		return nil, e
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
				oid, e := OidFromRef(repo, spec)
				if e != nil {
					return e
				}
				refs = append(refs, &ref{name: spec, oid: oid})
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
				oid, e := OidFromRef(repo, spec)
				if e != nil {
					return e
				}
				refs[spec] = &ref{name: spec, oid: oid}
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

func (repo *DiskRepository) RevParse(name string) (Object, error) {
	// TODO: implement
	return nil, nil
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

// validate a repository path to make sure it has
// the right format and that it exists
func validateRepo(path string) bool {
	// TODO
	return true
}
