//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import "github.com/jbrukh/ggit/api/objects"

// ================================================================= //
// FORMATTING
// ================================================================= //

// The "plumbing" string output of an Object. This output can be used
// to reproduce the contents or SHA1 hash of an Object.
func (f *Format) Object(o objects.Object) (int, error) {
	switch t := o.(type) {
	case *objects.Blob:
		return f.Blob(t)
	case *objects.Tree:
		return f.Tree(t)
	case *objects.Commit:
		return f.Commit(t)
	case *objects.Tag:
		return f.Tag(t)
	}
	panic("unknown object")
}

// The pretty string output of an Object. This format is not necessarily
// of use as an api call; it is for humans.
func (f *Format) ObjectPretty(o objects.Object) (int, error) {
	switch t := o.(type) {
	case *objects.Blob:
		return f.Blob(t)
	case *objects.Tree:
		return f.TreePretty(t)
	case *objects.Commit:
		return f.Commit(t)
	case *objects.Tag:
		return f.TagPretty(t)
	}
	panic("unknown object")
}

// ================================================================= //
// OPERATIONS
// ================================================================= //

// ObjectFromOid turns an ObjectId into an Object given the parent
// repository of the object.
func ObjectFromOid(repo Repository, oid *objects.ObjectId) (objects.Object, error) {
	return repo.ObjectFromOid(oid)
}

// ObjectFromOid turns a short hex into an Object given the parent
// repository of the object.
func ObjectFromShortOid(repo Repository, short string) (objects.Object, error) {
	return repo.ObjectFromShortOid(short)
}

// ObjectFromRef is similar to OidFromRef, except it derefernces the
// target ObjectId into an actual Object.
func ObjectFromRef(repo Repository, spec string) (objects.Object, error) {
	ref, err := PeeledRefFromSpec(repo, spec)
	if err != nil {
		return nil, err
	}
	return repo.ObjectFromOid(ref.ObjectId())
}

// ObjectFromRevision takes a revision specification and obtains the
// object that this revision specifies.
func ObjectFromRevision(repo Repository, rev string) (objects.Object, error) {
	p := newRevParser(repo, rev)
	e := p.Parse()
	if e != nil {
		return nil, e
	}
	return p.Object(), nil
}
