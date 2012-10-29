//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
benchmarks_test.go implements the ggit benchmarks.
*/
package api

import (
	"github.com/jbrukh/ggit/test"
	"testing"
)

func looseBlobOid() (repo Repository, oid *ObjectId) {
	testRepo := test.Blobs
	output := testRepo.Output().([]*test.OutputBlob)
	repo = Open(testRepo.Repo())
	oid = OidNow(output[0].Oid)
	return
}

func looseCommitOid() (repo Repository, oid *ObjectId) {
	testRepo := test.Linear
	output := testRepo.Output().(*test.OutputLinear)
	repo = Open(testRepo.Repo())
	oid = OidNow(output.Commits[0].Oid)
	return
}

func packedCommitOid() (repo Repository, oid *ObjectId) {
	testRepo := test.LinearPacked
	output := testRepo.Output().(*test.OutputLinearPacked)
	repo = Open(testRepo.Repo())
	oid = OidNow(output.Commits[0].Oid)
	return
}

func derefs() (repo Repository, output *test.OutputDerefs) {
	testRepo := test.Derefs
	output = testRepo.Output().(*test.OutputDerefs)
	repo = Open(testRepo.Repo())
	return
}

func Benchmark__readLooseBlobByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := looseBlobOid()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ObjectFromOid(oid)
		if err != nil {
			b.Fatalf("could not read object: %s", oid)
		}
	}
}

func Benchmark__readLooseCommitByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := looseCommitOid()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ObjectFromOid(oid)
		if err != nil {
			b.Fatalf("could not read object: %s", oid)
		}
	}
}

func Benchmark__readLooseTreeByOid(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	oid := OidNow(output.TreeOid)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ObjectFromOid(oid)
		if err != nil {
			b.Fatalf("could not read object: %s", oid)
		}
	}
}

func Benchmark__readLooseBlobByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := looseBlobOid()
	rev := oid.String()[:20]
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := ObjectFromRevision(repo, rev)
		if err != nil {
			b.Fatalf("could not read object: %s", rev)
		}
	}
}

func Benchmark__readLooseTreeByShort(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	rev := output.TreeOid[:20]
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := ObjectFromRevision(repo, rev)
		if err != nil {
			b.Fatalf("could not read object: %s", rev)
		}
	}
}

func Benchmark__readLooseCommitByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := looseCommitOid()
	rev := oid.String()[:20]
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := ObjectFromRevision(repo, rev)
		if err != nil {
			b.Fatalf("could not read object: %s", rev)
		}
	}
}

func Benchmark__readPackedCommitByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := packedCommitOid()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ObjectFromOid(oid)
		if err != nil {
			b.Fatalf("could not read object: %s", oid)
		}
	}
}

func Benchmark__readPackedTreeByOid(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	oid := OidNow(output.TreeOid)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ObjectFromOid(oid)
		if err != nil {
			b.Fatalf("could not read object: %s", oid)
		}
	}
}

func Benchmark__readPackedCommitByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := packedCommitOid()
	rev := oid.String()[:20]
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := ObjectFromRevision(repo, rev)
		if err != nil {
			b.Fatalf("could not read object: %s", rev)
		}
	}
}

func Benchmark__readPackedTreeByShort(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	rev := output.TreeOid[:20]
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := ObjectFromRevision(repo, rev)
		if err != nil {
			b.Fatalf("could not read object: %s", rev)
		}
	}
}

func Benchmark__derefLooseTreeFromCommit(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	oid := OidNow(output.CommitOid)
	rev := oid.String() + "^{tree}"
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := ObjectFromRevision(repo, rev)
		if err != nil {
			b.Fatalf("could not read object: %s", rev)
		}
	}
}
