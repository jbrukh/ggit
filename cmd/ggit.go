package main

import (
    "errors"
    "flag"
    "fmt"
    "github.com/jbrukh/ggit"
    "os"
)

func usage() {
    fmt.Println("USAGE: ggit <command> [<option>...] [<param>...]")
}

type handler func([]string) error

//
// FLAGS, flags everywhere. Put them in your car, put them in your
// wallet, put them in your cubicle, shove them up your ass.
//
var catFileFlags *flag.FlagSet = flag.NewFlagSet("cat-file", flag.ExitOnError)

var handlers map[string]handler = map[string]handler{
    "cat-file": catFile,
}

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        usage()
        os.Exit(1)
    }
    h, ok := handlers[args[0]]
    if !ok {
        usage()
        os.Exit(2)
    }
    if err := h(args); err != nil {
        fmt.Println("error: ", err)
        os.Exit(-1)
    }
    os.Exit(0)
}

func catFile(args []string) (err error) {
    // TODO: can we move this out?
    isType := catFileFlags.Bool("t", false, "show object type")
    isPrint := catFileFlags.Bool("p", false, "pretty-print object's contents")
    isSize := catFileFlags.Bool("s", false, "show object size")

    catFileFlags.Parse(args[1:])

    a := catFileFlags.Args()
    if len(a) != 1 {
        return errors.New("provide an object")
    }
    id := a[0]

    var (
        repo *ggit.Repository
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
