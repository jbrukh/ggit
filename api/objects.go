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
// OPERATIONS
// ================================================================= //

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
