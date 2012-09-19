package ggit

// ================================================================= //
// CONSTANTS RELATED TO TYPES
// ================================================================= //

// the types of Git objects
type ObjectType int8

const (
	OBJECT_BLOB ObjectType = iota
	OBJECT_TREE
	OBJECT_COMMIT
	OBJECT_TAG
)

// string representations of Git objects
const (
	OBJECT_BLOB_STR   = "blob"
	OBJECT_TREE_STR   = "tree"
	OBJECT_COMMIT_STR = "commit"
	OBJECT_TAG_STR    = "tag"
)

