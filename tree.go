package ggit

type FileMode uint16

const (
        MODE_FILE FileMode = 0100644
        MODE_EXEC FileMode = 0100755
        MODE_TREE FileMode = 0040000
        MODE_LINK FileMode = 0120000
)

type Tree struct {
        entries []TreeEntry
}

type TreeEntry struct {
        mode    FileMode
        otype   ObjectType
        name    string
        oid     *ObjectId
}
