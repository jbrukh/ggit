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
	output := testRepo.Output().(*test.OutputBlobs)
	repo = Open(testRepo.Repo())
	oid = OidNow(output.Blobs[0].Oid)
	return
}

func looseCommitOid() (repo Repository, oid *ObjectId) {
	testRepo := test.Linear
	output := testRepo.Output().(*test.OutputLinear)
	repo = Open(testRepo.Repo())
	oid = OidNow(output.Commits[0].Oid)
	return
}

func packedBlobOid() (repo Repository, oid *ObjectId) {
	testRepo := test.BlobsPacked
	output := testRepo.Output().(*test.OutputBlobsPacked)
	repo = Open(testRepo.Repo())
	oid = OidNow(output.Blobs[0].Oid)
	return
}

func packedCommitOid() (repo Repository, oid *ObjectId) {
	testRepo := test.LinearPacked
	output := testRepo.Output().(*test.OutputLinearPacked)
	repo = Open(testRepo.Repo())
	oid = OidNow(output.Commits[0].Oid)
	return
}

func looseDerefs() (repo Repository, output *test.OutputDerefs) {
	testRepo := test.Derefs
	output = testRepo.Output().(*test.OutputDerefs)
	repo = Open(testRepo.Repo())
	return
}

func packedDerefs() (repo Repository, output *test.OutputDerefs) {
	testRepo := test.DerefsPacked
	output = testRepo.Output().(*test.OutputDerefsPacked)
	repo = Open(testRepo.Repo())
	return
}

func objectFromOid(b *testing.B, repo Repository, oid *ObjectId) {
	b.StartTimer()
	_, err := repo.ObjectFromOid(oid)
	if err != nil {
		b.Fatalf("could not read object: %s", oid)
	}
	b.StopTimer()
}

func objectFromRev(b *testing.B, repo Repository, rev string) {
	b.StartTimer()
	_, err := ObjectFromRevision(repo, rev)
	if err != nil {
		b.Fatalf("could not read object: %s", rev)
	}
	b.StopTimer()
}

func Benchmark__readLooseBlobByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := looseBlobOid()
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readLooseCommitByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := looseCommitOid()
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readLooseTreeByOid(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	oid := OidNow(output.TreeOid)
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readLooseBlobByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := looseBlobOid()
	rev := oid.String()[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__readLooseTreeByShort(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	rev := output.TreeOid[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__readLooseCommitByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := looseCommitOid()
	rev := oid.String()[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__readPackedBlobByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := packedBlobOid()
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readPackedCommitByOid(b *testing.B) {
	b.StopTimer()
	repo, oid := packedCommitOid()
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readPackedTreeByOid(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	oid := OidNow(output.TreeOid)
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readPackedCommitByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := packedCommitOid()
	rev := oid.String()[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__readPackedTreeByShort(b *testing.B) {
	b.StopTimer()
	repo, output := derefs()
	rev := output.TreeOid[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__derefLooseTreeFromCommit(b *testing.B) {
	b.StopTimer()
	repo, output := looseDerefs()
	oid := OidNow(output.CommitOid)
	rev := oid.String() + "^{tree}"
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__derefPackedTreeFromCommit(b *testing.B) {
	b.StopTimer()
	repo, output := packedDerefs()
	oid := OidNow(output.CommitOid)
	rev := oid.String() + "^{tree}"
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}
