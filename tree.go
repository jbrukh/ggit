package ggit

type FileMode uint16

const (
    MODE_FILE FileMode = 0100644
    MODE_EXEC FileMode = 0100755
    MODE_TREE FileMode = 0040000
    MODE_LINK FileMode = 0120000
)

type rawTree struct {
	rawObj *RawObject
}

func (rt *rawTree) Parse() (err error) {
	return
}

func (rt *rawTree) ParseTree() (t *Tree, err error) {
	return nil, nil
}

func newRawTree(rawObj *RawObject) (rt *rawTree) {
	return &rawTree{
		rawObj,
	}
}

/*func parseModeAndName(b []byte) (modeStr, nameStr string, hInx uint) {
    const MAX_SZ = 32
    var i, j uint
    for i = 0; i < MAX_SZ; i++ {
        if b[i] == ' ' {
            typeStr = string(b[:i])
            for j = i; j < MAX_SZ; j++ {
                if b[j] == '\000' {
                    pInx = j
                    sizeStr = string(b[i+1 : j])
                    return
                }
            }
        }
    }
    return
}*/
	
	/*b := bytes.NewBuffer(p)
    for {
        modeName, err := b.ReadString('\000')
        if err != nil {
            break
        }
        fmt.Printf("%v\n", modeName)
        bts := b.Next(20)
        hsh := NewObjectIdFromBytes(bts)
        fmt.Printf("sha: %s\n", hsh.String())
        if err != nil {
            break
        }
    }*/

type Tree struct {
    entries []TreeEntry
	parent *Repository
}

type TreeEntry struct {
    mode  FileMode
    otype ObjectType
    name  string
    oid   *ObjectId
}
