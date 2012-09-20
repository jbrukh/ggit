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
    ObjectBlob.String(),
    ObjectTree.String(),
    ObjectCommit.String(),
    ObjectTag.String(),
}
