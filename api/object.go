package api

type Object interface {
	// Type returns the type of this object. Available types are 
	// defined by ObjectType are usually one of blob, tree, 
	// commit, or tag.
	Type() ObjectType

	// Size returns the size of the payload of this object.
	Size() int

	// Some default representation of this object.
	String() string
}
