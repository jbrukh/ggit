package ggit

import "errors"

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

// return a human-readable representation of an ObjectType
// TODO: turn this into a to-function
func (otype ObjectType) String() string {
	switch otype {
	case OBJECT_BLOB:
		return OBJECT_BLOB_STR
	case OBJECT_TREE:
		return OBJECT_TREE_STR
	case OBJECT_COMMIT:
		return OBJECT_COMMIT_STR
	case OBJECT_TAG:
		return OBJECT_TAG_STR
	}
	panic("unknown type")
}

func toObjectType(typeStr string) (otype ObjectType, err error) {
	switch typeStr {
	case OBJECT_BLOB_STR:
		return OBJECT_BLOB, nil
	case OBJECT_TREE_STR:
		return OBJECT_TREE, nil
	case OBJECT_TAG_STR:
		return OBJECT_TAG, nil
	case OBJECT_COMMIT_STR:
		return OBJECT_COMMIT, nil
	}
	return 0, errors.New("unknown object type")
}
