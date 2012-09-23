package api

import (
    "errors"
    "flag"
    "fmt"
    "github.com/jbrukh/ggit"
    "os"
)

//
// FLAGS, flags everywhere. Put them in your car, put them in your
// wallet, put them in your cubicle, shove them up your ass.
//
var catFileFlags *flag.FlagSet = flag.NewFlagSet("cat-file", flag.ExitOnError)

var (
    isType  *bool
    isPrint *bool
    isSize  *bool
)

func init() {
    isType = catFileFlags.Bool("t", false, "show object type")
    isPrint = catFileFlags.Bool("p", false, "pretty-print object's contents")
    isSize = catFileFlags.Bool("s", false, "show object size")
}

func CatFile(args []string) (err error) {

    catFileFlags.Parse(args[1:])

    a := catFileFlags.Args()
    if len(a) < 1 {
        return errors.New("provide an object")
    }
    if len(a) > 1 {
        return errors.New(fmt.Sprint("expecting a single argument, found ", len(a)))
    }
    id := a[0]

    var (
        repo *ggit.DiskRepository
        oid  *ggit.ObjectId
    )

    // get a proper id
    if oid, err = ggit.NewObjectIdFromString(id); err != nil {
        return
    }

    // TODO: perhaps not open the repo before parsing args?
    if repo, err = ggit.Open(ggit.DEFAULT_GIT_DIR); err != nil {
        return
    }
    defer repo.Close()

    switch {
    case *isPrint:
        return doPrint(repo, oid)
    case *isType:
        return doType(repo, oid)
    case *isSize:
        return doSize(repo, oid)
    default:
        return errors.New("unknown command")
    }
    return
}

func doPrint(repo *ggit.DiskRepository, oid *ggit.ObjectId) (err error) {
    if obj, err := repo.ReadObject(oid); err == nil {
        obj.WriteTo(os.Stdout)
        return err
    }
    return errors.New("could not find object: " + oid.String() + ": " + err.Error()) // TODO
}

func doType(repo *ggit.DiskRepository, oid *ggit.ObjectId) (err error) {
    var obj ggit.Object
    if obj, err = repo.ReadObject(oid); err != nil {
        return err
    }
    fmt.Println(obj.Type())
    return
}

// commenting until I figure out what size means in this context
func doSize(repo *ggit.DiskRepository, oid *ggit.ObjectId) (err error) {
    var obj ggit.Object
    if obj, err = repo.ReadObject(oid); err != nil {
        return err
    }
    fmt.Println(obj.Size())
    return
}
