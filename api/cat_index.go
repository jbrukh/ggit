package api

import (
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
    var (
        repo *ggit.Repository
    )

    // TODO: perhaps not open the repo before parsing args?
    if repo, err = ggit.Open(ggit.DEFAULT_GIT_DIR); err != nil {
        return
    }
    defer repo.Close()

    ggit.ParseIndexFile(repo)

    return
}
