package api

type Object interface {
	// Type returns the type of this object. Available types are 
	// defined by ObjectType are usually one of blob, tree, 
	// commit, or tag.
	Type() ObjectType

	// Size returns the size of the payload of this object.
	Size() int
}

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
