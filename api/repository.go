//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import ()

const (
	DefaultGitDir     = ".git"
	DefaultObjectsDir = "objects"
	DefaultPackDir    = "pack"
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
