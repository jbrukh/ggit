package ggit

type Tree struct {
    parent Tree
    object interface{} // either a tree or a blob
    children []TreeEntry
}

type (t *Tree) Type() ObjectType {
    return OBJECT_TREE
}

func (t *Tree) WriteTo(io.Writer) (id *ObjectId, err error) {
    // TODO: figure out the format of the tree
    return nil, nil
}
