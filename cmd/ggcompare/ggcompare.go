//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package main

import (
	"bytes"
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

	oids, e := api.ObjectIds(repo)
	if e != nil {
		fmt.Fprintf(os.Stderr, "could not open repo: %s\n", e.Error())
		os.Exit(1)
	}

	for i, oid := range oids {
		fmt.Println(i, "git/ggit object:", oid)

		o, err := repo.ObjectFromOid(oid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ggit could not read object: %s\n", oid)
			os.Exit(1)
		}

		// get the type for git cat-file
		t := o.Header().Type()
		// get the git cat-file output
		var cat string
		cat, err = util.GitExec(repoPath, "cat-file", t.String(), oid.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "git could not read object: %s\n", oid)
			os.Exit(1)
		}

		f := api.NewStrFormat()
		f.Object(o)
		str := f.String()
		if str != cat {
			fmt.Printf("\nfound a mismatch in object %d (%s)...\n", i, oid)
			fmt.Println("GIT -----------------------")
			fmt.Println(cat)
			fmt.Println("GGIT -----------------------")
			fmt.Println(str)
			fmt.Println("diff -----------------------")

			i := 0
			for {
				if cat[i] != str[i] {
					break
				}
				i++
			}

			ours := bytes.NewBufferString(str)
			theirs := bytes.NewBufferString(cat)
			b := ours.Next(i)
			theirs.Next(i)

			fmt.Println(string(b))
			line1, _ := ours.ReadString('\n')
			line2, _ := theirs.ReadString('\n')

			fmt.Printf("theirs:\t%s", line2)
			fmt.Printf("ours:\t%s", line1)

			os.Exit(2)
		}
	}
}
