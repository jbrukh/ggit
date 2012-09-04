package ggit

import (
    "io"
)

type Tree struct {
    // TODO
}

func (t *Tree) Type() ObjectType {
    return OBJECT_TREE
}

func (t *Tree) WriteTo(io.Writer) (id *ObjectId, err error) {
    // TODO: figure out the format of the tree
    return nil, nil
}
