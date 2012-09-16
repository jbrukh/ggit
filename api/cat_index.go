package api

import (
    "fmt"
    "github.com/jbrukh/ggit"
)

//
// FLAGS, flags everywhere. Put them in your car, put them in your
// wallet, put them in your cubicle, shove them up your ass.
//
/*
var catIndexFlags *flag.FlagSet = flag.NewFlagSet("cat-index", flag.ExitOnError)

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
*/

func CatIndex(args []string) (err error) {

    repo, e := ggit.Open(ggit.DEFAULT_GIT_DIR)
    if e != nil {
        return e
    }
    defer repo.Close()

    inx, e := repo.Index()
    if e != nil {
        return e
    }
    fmt.Print(inx)

    return
}
