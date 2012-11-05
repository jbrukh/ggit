//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"fmt"
)

type RevVisitor func(c *Commit) error

var RevPrinter RevVisitor = func(c *Commit) error {
	fmt.Println(c.ObjectId())
	return nil
}

func RevWalkFromRevision(repo Repository, rev string, f RevVisitor) error {
	root, err := CommitFromRevision(repo, rev)
	if err != nil {
		return err
	}

	current := root
	for {
		// apply the visitor
		if err = f(current); err != nil {
			return err
		}

		// proceed to the next commit
		parents := current.Parents()
		if len(parents) < 1 {
			break
		}
		current, err = CommitNthParent(repo, current, 1)
		if err != nil {
			return err
		}
	}
	return err
}
