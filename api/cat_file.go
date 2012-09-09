package api

import (
	"errors"
	"flag"
	"os"
	"fmt"
	. "github.com/jbrukh/ggit"
)

//
// FLAGS, flags everywhere. Put them in your car, put them in your
// wallet, put them in your cubicle, shove them up your ass.
//
var catFileFlags *flag.FlagSet = flag.NewFlagSet("cat-file", flag.ExitOnError)

var (
	isType *bool 
	isPrint *bool 
	isSize *bool
)

func init() {
	isType = catFileFlags.Bool("t", false, "show object type")
	isPrint = catFileFlags.Bool("p", false, "pretty-print object's contents")
	isSize = catFileFlags.Bool("s", false, "show object size")
}

func CatFile(args []string) (err error) {

    catFileFlags.Parse(args[1:])

    a := catFileFlags.Args()
    if len(a) != 1 {
        return errors.New("provide an object")
    }
    id := a[0]

    var (
        repo *Repository
        oid  *ObjectId
    )

    // get a proper id
    if oid, err = NewObjectIdFromString(id); err != nil {
        return
    }

    // TODO: perhaps not open the repo before parsing args?
    if repo, err = Open(DEFAULT_GIT_DIR); err != nil {
        return
    }
    defer repo.Close()

    if *isPrint {
        obj, e := repo.ReadObject(oid)
        if e != nil {
            return errors.New("could not find object: " + oid.String()) // TODO
        }
        obj.WriteTo(os.Stdout)
    } else if *isType {
        h, e := repo.ReadRawObjectHeader(oid)
        if e != nil {
            return e
        }
        fmt.Println(h.Type)
    } else if *isSize {
        h, e := repo.ReadRawObjectHeader(oid)
        if e != nil {
            return e
        }
        fmt.Println(h.Size)
    } else {
        catFileFlags.PrintDefaults()
    }

    return
}