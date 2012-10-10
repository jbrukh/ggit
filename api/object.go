package api

// ObjectHeader contains the type and size
// information for an object.
type ObjectHeader interface {
	Type() ObjectType
	Size() int
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
