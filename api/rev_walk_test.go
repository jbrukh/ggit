package api

import (
	"container/heap"
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_PriorityQueue(t *testing.T) {
	mockCommit := func(name string, seconds int64) *Commit {
		return &Commit{
			oid: OidNow(name),
			committer: &WhoWhen{
				When: When{
					seconds: seconds,
				},
			},
		}
	}
	commits := []*Commit{
		mockCommit("0000000000000000000000000000000000000000", 0),
		mockCommit("1111111111111111111111111111111111111111", 1),
		mockCommit("2222222222222222222222222222222222222222", 2),
		mockCommit("3333333333333333333333333333333333333333", 3),
		mockCommit("4444444444444444444444444444444444444444", 4),
		mockCommit("5555555555555555555555555555555555555555", 5),
	}

	pq := PriorityQueue(make([]*Commit, 0))
	util.Assert(t, len(pq) == 0)

	heap.Init(&pq)
	heap.Push(&pq, commits[5])
	heap.Push(&pq, commits[0])
	heap.Push(&pq, commits[4])
	heap.Push(&pq, commits[3])
	heap.Push(&pq, commits[1])
	heap.Push(&pq, commits[2])

	util.Assert(t, len(pq) == len(commits))
	next := func() int {
		return int(heap.Pop(&pq).(*Commit).Committer().Seconds())
	}
	util.AssertEqualInt(t, 5, next())
	util.AssertEqualInt(t, 4, next())
	util.AssertEqualInt(t, 3, next())
	util.AssertEqualInt(t, 2, next())
	util.AssertEqualInt(t, 1, next())
	util.AssertEqualInt(t, 0, next())
}
