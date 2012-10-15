//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package main

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/util"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: ggcompare <repo>\n")
		os.Exit(1)
	}

	repoPath := args[0]
	println(repoPath)

	repo := api.Open(repoPath)

	oids, e := repo.ObjectIds()
	if e != nil {
		fmt.Fprintf(os.Stderr, "could not open repo: %s\n", e.Error())
		os.Exit(1)
	}

	for _, oid := range oids {
		fmt.Println(oid)

		// get the dashP from git
		dashP, err := util.GitExec(repoPath, "cat-file", "-p", oid.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "git could not read object: %s\n", oid)
			os.Exit(1)
		}

		var o api.Object
		o, err = repo.ObjectFromOid(oid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ggit could not read object: %s\n", oid)
			os.Exit(1)
		}

		f := api.NewStrFormat()
		f.Object(o)
		str := f.String()
		if str != dashP {
			fmt.Printf("mismatch; expected '%d' but got '%d'\n", len(dashP), len(str))
			i := 0
			for {
				if dashP[i] != str[i] {
					break
				}
				i++
			}
			fmt.Println(dashP[:i])
			fmt.Println(dashP[i:i+10], "\n", str[i:i+10])
			fmt.Println("-=====---==-==")
			fmt.Println(str)
			os.Exit(2)
		}
	}
}
