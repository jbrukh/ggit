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

// ================================================================= //
// OBJECT PROVIDERS
// ================================================================= //

func looseBlobOid() (repo Repository, oid *ObjectId) {
	testRepo := test.Derefs
	info := testRepo.Info().(*test.InfoDerefs)
	repo = Open(testRepo.Repo())
	oid = OidNow(info.BlobOid)
	return
}

func looseCommitOid() (repo Repository, oid *ObjectId) {
	testRepo := test.Linear
	info := testRepo.Info().(*test.InfoLinear)
	repo = Open(testRepo.Repo())
	oid = OidNow(info.Commits[1].CommitOid)
	return
}

func packedBlobOid() (repo Repository, oid *ObjectId) {
	testRepo := test.DerefsPacked
	info := testRepo.Info().(*test.InfoDerefsPacked)
	repo = Open(testRepo.Repo())
	oid = OidNow(info.BlobOid)
	return
}

func packedCommitOid() (repo Repository, oid *ObjectId) {
	testRepo := test.LinearPacked
	info := testRepo.Info().(*test.InfoLinearPacked)
	repo = Open(testRepo.Repo())
	oid = OidNow(info.Commits[1].CommitOid)
	return
}

func looseDerefs() (repo Repository, info *test.InfoDerefs) {
	testRepo := test.Derefs
	info = testRepo.Info().(*test.InfoDerefs)
	repo = Open(testRepo.Repo())
	return
}

func packedDerefs() (repo Repository, info *test.InfoDerefsPacked) {
	testRepo := test.DerefsPacked
	info = testRepo.Info().(*test.InfoDerefsPacked)
	repo = Open(testRepo.Repo())
	return
}

func looseRef() (repo Repository, spec string) {
	testRepo := test.Refs
	repo = Open(testRepo.Repo())
	info := testRepo.Info().(*test.InfoRefs)
	spec = info.BranchName
	return
}

func packedRef() (repo Repository, spec string) {
	testRepo := test.Refs
	repo = Open(testRepo.Repo())
	info := testRepo.Info().(*test.InfoRefs)
	spec = info.AnnTagName
	return
}

func justBlob() (Repository, *Blob) {
	repo, oid := packedBlobOid()
	o, err := repo.ObjectFromOid(oid)
	if err != nil {
		panic(err)
	}
	return repo, o.(*Blob)
}

func justCommit() (Repository, *Commit) {
	repo, oid := packedCommitOid()
	o, err := repo.ObjectFromOid(oid)
	if err != nil {
		panic(err)
	}
	return repo, o.(*Commit)
}

func justTree() (Repository, *Tree) {
	repo, commit := justCommit()
	o, err := repo.ObjectFromOid(commit.Tree())
	if err != nil {
		panic(err)
	}
	return repo, o.(*Tree)
}

func justTag() (Repository, *Tag) {
	repo, info := packedDerefs()
	tagOid := OidNow(info.TagOid)
	o, err := repo.ObjectFromOid(tagOid)
	if err != nil {
		panic(err)
	}
	return repo, o.(*Tag)
}

// ================================================================= //
// ATOM OPERATIONS
// ================================================================= //

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

func listRefs(b *testing.B, repo Repository) {
	b.StartTimer()
	_, err := repo.Refs()
	if err != nil {
		b.Fatalf("could not list refs: %s", err)
	}
	b.StopTimer()
}

func resolveRef(b *testing.B, repo Repository, spec string) {
	b.StartTimer()
	_, err := repo.Ref(spec)
	if err != nil {
		b.Fatalf("could not resolve ref: %s", err)
	}
	b.StopTimer()
}

func resolveShortRef(b *testing.B, repo Repository, spec string) {
	b.StartTimer()
	_, err := RefFromSpec(repo, spec)
	if err != nil {
		b.Fatalf("could not resolve ref: %s", err)
	}
	b.StopTimer()
}

// ================================================================= //
// BENCHMARKS - FORMATTING
// ================================================================= //

func Benchmark__oidFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OidNow("fb5b685ca75023e129b2f3c8172a7a166ad4dca8")
	}
}

func Benchmark__oidToString(b *testing.B) {
	b.StopTimer()
	const oidStr = "fb5b685ca75023e129b2f3c8172a7a166ad4dca8"
	for i := 0; i < b.N; i++ {
		oid := OidNow(oidStr)
		b.StartTimer()
		oid.String()
		b.StopTimer()
	}
}

// ================================================================= //
// BENCHMARKS - FORMATTING
// ================================================================= //

func Benchmark__blobString(b *testing.B) {
	b.StopTimer()
	f := NewStrFormat()
	_, blob := justBlob()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f.Blob(blob)
	}
}

func Benchmark__commitString(b *testing.B) {
	b.StopTimer()
	f := NewStrFormat()
	_, commit := justCommit()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f.Commit(commit)
	}
}

func Benchmark__treeString(b *testing.B) {
	b.StopTimer()
	f := NewStrFormat()
	_, tree := justTree()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f.Tree(tree)
	}
}

func Benchmark__tagString(b *testing.B) {
	b.StopTimer()
	f := NewStrFormat()
	_, tag := justTag()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f.Tag(tag)
	}
}

// ================================================================= //
// BENCHMARKS - REFS
// ================================================================= //

func Benchmark__listTenRefs(b *testing.B) {
	b.StopTimer()
	repo := Open(test.Linear.Repo())
	for i := 0; i < b.N; i++ {
		listRefs(b, repo)
	}
}

func Benchmark__resolveLooseRef(b *testing.B) {
	b.StopTimer()
	repo, spec := looseRef()
	for i := 0; i < b.N; i++ {
		resolveRef(b, repo, expandHeadRef(spec)) // TODO
	}
}

func Benchmark__resolvePackedRef(b *testing.B) {
	b.StopTimer()
	repo, spec := packedRef()
	for i := 0; i < b.N; i++ {
		resolveRef(b, repo, expandTagRef(spec)) // TODO
	}
}

func Benchmark__resolveLooseShortRef(b *testing.B) {
	b.StopTimer()
	repo, spec := looseRef()
	for i := 0; i < b.N; i++ {
		resolveShortRef(b, repo, spec)
	}
}

func Benchmark__resolvePackedShortRef(b *testing.B) {
	b.StopTimer()
	repo, spec := packedRef()
	for i := 0; i < b.N; i++ {
		resolveShortRef(b, repo, spec)
	}
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
	repo, info := looseDerefs()
	oid := OidNow(info.TreeOid)
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

func Benchmark__readLooseCommitByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := looseCommitOid()
	rev := oid.String()[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__readLooseTreeByShort(b *testing.B) {
	b.StopTimer()
	repo, info := looseDerefs()
	rev := info.TreeOid[:20]
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
	repo, info := packedDerefs()
	oid := OidNow(info.TreeOid)
	for i := 0; i < b.N; i++ {
		objectFromOid(b, repo, oid)
	}
}

func Benchmark__readPackedBlobByShort(b *testing.B) {
	b.StopTimer()
	repo, oid := packedBlobOid()
	rev := oid.String()[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
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
	repo, info := packedDerefs()
	rev := info.TreeOid[:20]
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__derefLooseTreeFromCommit(b *testing.B) {
	b.StopTimer()
	repo, info := looseDerefs()
	oid := OidNow(info.CommitOid)
	rev := oid.String() + "^{tree}"
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__derefLooseTreeFromBranch(b *testing.B) {
	b.StopTimer()
	repo, info := looseDerefs()
	rev := info.BranchName + "^{tree}"
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__derefPackedTreeFromCommit(b *testing.B) {
	b.StopTimer()
	repo, info := packedDerefs()
	oid := OidNow(info.CommitOid)
	rev := oid.String() + "^{tree}"
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__derefPackedTreeFromBranch(b *testing.B) {
	b.StopTimer()
	repo, info := packedDerefs()
	rev := info.BranchName + "^{tree}"
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__resolveLooseBranch(b *testing.B) {
	b.StopTimer()
	repo, info := looseDerefs()
	rev := info.BranchName
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}

func Benchmark__resolvePackedBranch(b *testing.B) {
	b.StopTimer()
	repo, info := packedDerefs()
	rev := info.BranchName
	for i := 0; i < b.N; i++ {
		objectFromRev(b, repo, rev)
	}
}
