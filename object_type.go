package ggit

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

func assertObjectType(str string) (ObjectType, bool) {
	otype := ObjectType(str)
	switch otype {
	case ObjectBlob, ObjectCommit, ObjectTag, ObjectTree:
		return otype, true
	}
	return 0, false
}
