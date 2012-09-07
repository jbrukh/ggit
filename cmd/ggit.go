package main

import (
        "errors"
        "flag"
        "fmt"
        "github.com/jbrukh/ggit"
        "os"
)

func usage() {
        fmt.Println("USAGE: ggit <command>")
}

type handler func([]string) error

var handlers map[string]handler = map[string]handler{
        "hash-object": hashObject,
        "read-blob":   readBlob,
        "read-tree":   readTree,
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
                fmt.Println(err)
        }
}

func hashObject(args []string) error {
        return nil
}

func readBlob(args []string) (err error) {
        if len(args) < 2 {
                return errors.New("provide a hash")
        }
        oid, err := ggit.NewObjectIdFromString(args[1])
        repo, _ := ggit.Open(".git")
        b, err := repo.ReadBlob(oid)
        if err != nil {
                fmt.Println("Error: ", err)
        }
        fmt.Println("bytes: ", string(b.Bytes()))
        return
}

func readTree(args []string) (err error) {
        if len(args) < 2 {
                return errors.New("provide a hash")
        }
        oid, err := ggit.NewObjectIdFromString(args[1])
        repo, _ := ggit.Open(".git")
        err = repo.ReadTree(oid)
        if err != nil {
                fmt.Println("Error: ", err)
        }
        return
}
