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
    "cat-file":  catFile,
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
	fs := flag.NewFlagSet("cat-file", flag.ExitOnError)
	oidForType := fs.String("t", "", "show object type")
	oidForPrint := fs.String("p", "", "pretty-print object's contents")
	fs.Parse(args[1:])
	
	if *oidForPrint != "" {
		repo, e := ggit.Open(".git")
		if e != nil {
			return e
		}
		defer repo.Close()
		oid, e := ggit.NewObjectIdFromString(*oidForPrint)
		if e != nil {
			return e
		}
		
		obj, e := repo.ReadObject(oid)
		if e != nil {
			return errors.New("could not find object: " + oid.String()) // TODO
		}
		
		obj.WriteTo(os.Stdout)
	} else if *oidForType != "" {
		repo, e := ggit.Open(".git")
		if e != nil {
			return e
		}
		defer repo.Close()
		oid, e := ggit.NewObjectIdFromString(*oidForType)
		if e != nil {
			return e
		}
		
		obj, e := repo.ReadRawObject(oid)
		if e != nil {
			return errors.New("could not find object: " + oid.String()) // TODO
		}
		
		h, e := obj.Header()
		if e != nil {
			return e
		}

		fmt.Println(h.Type)
	} else {
		fs.PrintDefaults()
	}
	return
}

