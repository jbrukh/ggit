//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package objects

// ================================================================= //
// REF OBJECTS
// ================================================================= //

// Ref is a representation of a ggit reference. A ref is a nice
// name for an ObjectId. More precisely, a ref is a path relative
// to the git directory (without duplicate path separators, ".", or
// "..").
type Ref interface {

	// Name returns the string name of this ref. This is
	// a simple path relative to the git directory, which
	// may or may not be HEAD, MERGE_HEAD, etc.
	Name() string

	// Target returns the target reference, whether an oid
	// or another string ref. If the ref is symbolic then
	// "symbolic" is true.
	Target() (symbolic bool, o interface{})

	// ObjectId returns the object id that this ref references
	// provided this ref is not symbolic, and otherwise panics.
	ObjectId() *ObjectId

	// If this ref is a tag, then this field may contain
	// the target commit of the tag, if such an optimization
	// is available. Otherwise, this field is nil.
	Commit() *ObjectId
}

// ================================================================= //
// REF IMPLEMENTATION
// ================================================================= //

type ref struct {
	name   string
	oid    *ObjectId
	spec   string
	commit *ObjectId // if tag, this is the commit the tag points to
}

func NewRef(name, spec string, oid, commit *ObjectId) Ref {
	return &ref{name, oid, spec, commit}
}

func (r *ref) Target() (bool, interface{}) {
	if r.oid != nil {
		return false, r.oid
	}
	if r.spec != "" {
		return true, r.spec
	}
	panic("does not have an object reference")
}

func (r *ref) Name() string {
	return r.name
}

func (r *ref) Commit() *ObjectId {
	return r.commit
}

func (r *ref) ObjectId() *ObjectId {
	symbolic, oid := r.Target()
	if !symbolic {
		return oid.(*ObjectId)
	}
	panic("cannot return oid: this is a symbolic ref")
}
