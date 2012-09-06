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
var handlers map[string]handler = map[string]handler {
	"hash-object": hashObject,
    "cat-file": catFile,
	"read-blob": readBlob,
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

func catFile(args []string) (err error) {
    if len(args) < 2 {
        return errors.New("provide a file to cat-file")
    }
    file := args[1]
    b, err := ggit.NewBlobFromFile(file)
    if err == nil {
        fmt.Println(b.String())
    } else {
        fmt.Println("error: ", err)
    }
    return
}

func readBlob(args []string) (err error) {
	if len(args) < 2 {
		return errors.New("provide a hash")
	}
	oid, err := ggit.NewObjectIdFromString(args[1])
	repo, _ := ggit.OpenRepository(".git")
	b, err := repo.ReadBlob(oid)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(b)
	return
}
