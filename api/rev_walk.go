//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"container/heap"
	"fmt"
)

type RevVisitor func(c *Commit) error

var RevPrinter RevVisitor = func(c *Commit) error {
	fmt.Printf("%s %s\n", c.ObjectId(), c.Committer().Seconds())
	return nil
}

type RevWalkOptions struct {
	NoMerges  bool
	DateOrder bool
}

func RevWalkDateOrder(repo Repository, rev string, opts *RevWalkOptions, f RevVisitor) error {
	root, err := CommitFromRevision(repo, rev)
	if err != nil {
		return err
	}

	pq := PriorityQueue(make([]*Commit, 0))
	seen := make(map[string]bool, 0)
	seen[root.ObjectId().String()] = true
	heap.Push(&pq, root)
	var curr *Commit

	for len(pq) > 0 {
		// get the commit
		curr = heap.Pop(&pq).(*Commit)

		// apply the visitor
		if err = f(curr); err != nil {
			return err
		}

		// put the parents in there
		parents := curr.Parents()
		for _, oid := range parents {
			key := oid.String()
			if ok := seen[key]; ok {
				continue
			}
			p, err := CommitFromOid(repo, oid)
			if err != nil {
				return err
			}
			seen[key] = true
			heap.Push(&pq, p)
		}
	}
	return nil
}

func RevWalkFromRevision(repo Repository, rev string, opts *RevWalkOptions, f RevVisitor) error {
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

// ================================================================= //
// PRIORITY QUEUE
// ================================================================= //

type PriorityQueue []*Commit

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Committer().Seconds() > pq[j].Committer().Seconds()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Commit))
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	if n < 1 {
		return nil
	}
	result := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return result
}
