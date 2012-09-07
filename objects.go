package ggit

import (
)

type ObjectType int

// the types of objects
const (
    OBJECT_BLOB ObjectType = iota
    OBJECT_TREE
    OBJECT_COMMIT
    OBJECT_TAG
)
