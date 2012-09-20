package ggit

import (
	"errors"
	"fmt"
)

// ================================================================= //
// CONSTANTS RELATED TO TYPES
// ================================================================= //

// the types of Git objects
type ObjectType string

// return a human-readable representation of an ObjectType
func (otype ObjectType) String() string {
    return string(otype)
}

const (
    ObjectBlob   ObjectType = "blob"
    ObjectTree   ObjectType = "tree"
    ObjectCommit ObjectType = "commit"
    ObjectTag    ObjectType = "tag"
)

var objectTypes []string = []string{
    ObjectBlob.String(),
    ObjectTree.String(),
    ObjectCommit.String(),
    ObjectTag.String(),
}

func toObjectType(s string) (t ObjectType, err error) {
	switch s {
	case ObjectBlob.String():
		t = ObjectBlob
	case ObjectTree.String():
		t = ObjectTree
	case ObjectCommit.String():
		t = ObjectCommit
	case ObjectTag.String():
		t = ObjectTag
	default:
		err = errors.New(fmt.Sprint("Unrecognized object type. Expected one of: ", objectTypes))
	}
	return
}
