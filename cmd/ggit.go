package main

import (
    "flag"
    "fmt"
    "github.com/jbrukh/ggit/api"
    "os"
)

func usage() {
    fmt.Println("USAGE: ggit <command> [<option>...] [<param>...]")
}

type handler func([]string) error

var handlers map[string]handler = map[string]handler{
    "cat-file": api.CatFile,
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
