//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

// ObjectHeader contains the type and size
// information for an object.
type ObjectHeader interface {
	Type() ObjectType
	Size() int64
}

// Object represents a generic git object: a blob, a tree,
// a tag, or a commit.
type Object interface {

	// Header returns the object header, which
	// contains the object's type and size.
	Header() ObjectHeader

	// ObjectId returns the object id of the object.
	ObjectId() *ObjectId
}

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) Object(o Object) (int, error) {
	switch t := o.(type) {
	case *Blob:
		return f.Blob(t)
	case *Tree:
		return f.Tree(t)
	case *Commit:
		return f.Commit(t)
	case *Tag:
		return f.Tag(t)
	}
	panic("unknown object")
}

// ================================================================= //
// OPERATIONS
// ================================================================= //

// ObjectFromOid turns an ObjectId into an Object given the parent
// repository of the object.
func ObjectFromOid(repo Repository, oid *ObjectId) (Object, error) {
	return repo.ObjectFromOid(oid)
}

// ObjectFromOid turns a short hex into an Object given the parent
// repository of the object.
func ObjectFromShortOid(repo Repository, short string) (Object, error) {
	return repo.ObjectFromShortOid(short)
}

// ObjectFromRef is similar to OidFromRef, except it derefernces the
// target ObjectId into an actual Object.
func ObjectFromRef(repo Repository, spec string) (Object, error) {
	ref, err := PeeledRefFromSpec(repo, spec)
	if err != nil {
		return nil, err
	}
	return repo.ObjectFromOid(ref.ObjectId())
}

// ObjectFromRevision takes a revision specification and obtains the
// object that this revision specifies.
func ObjectFromRevision(repo Repository, rev string) (Object, error) {
	p := newRevParser(repo, rev)
	e := p.Parse()
	if e != nil {
		return nil, e
	}
	return p.Object(), nil
}
