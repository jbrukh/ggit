//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

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
	string(ObjectBlob),
	string(ObjectTree),
	string(ObjectCommit),
	string(ObjectTag),
}

// Return the ObjectType this string represents.
// If the string is invalid, return ObjectType(0), false
func AssertObjectType(str string) (ObjectType, bool) {
	otype := ObjectType(str)
	switch otype {
	case ObjectBlob, ObjectCommit, ObjectTag, ObjectTree:
		return otype, true
	}
	return ObjectType(0), false
}
