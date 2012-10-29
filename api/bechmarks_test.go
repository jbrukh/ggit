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

func Benchmark__readLooseBlobByOid(b *testing.B) {
	b.StopTimer()
	testRepo := test.Blobs
	repo := Open(testRepo.Repo())
	output := testRepo.Output().([]*test.OutputBlob)
	oid := OidNow(output[0].Oid)

	// make sure this operation is valid
	_, err := repo.ObjectFromOid(oid)
	if err != nil {
		b.Errorf("could not read blob: %s", oid)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		repo.ObjectFromOid(oid)
	}
}

func Benchmark__readLooseBlobByShort(b *testing.B) {
	b.StopTimer()
	testRepo := test.Blobs
	repo := Open(testRepo.Repo())
	output := testRepo.Output().([]*test.OutputBlob)
	oid := output[0].Oid[:20]

	// make sure this operation is valid
	_, err := ObjectFromRevision(repo, oid)
	if err != nil {
		b.Errorf("could not read blob: %s", oid)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ObjectFromRevision(repo, oid)
	}
}

func Benchmark__readLooseCommitByOid(b *testing.B) {
	b.StopTimer()
	testRepo := test.Linear
	repo := Open(testRepo.Repo())
	output := testRepo.Output().(*test.OutputLinear)
	oid := OidNow(output.Commits[0].Oid)

	// make sure this operation is valid
	_, err := repo.ObjectFromOid(oid)
	if err != nil {
		b.Errorf("could not read blob: %s", oid)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		repo.ObjectFromOid(oid)
	}
}

func Benchmark__readPackedCommitByOid(b *testing.B) {
	b.StopTimer()
	testRepo := test.LinearPacked
	repo := Open(testRepo.Repo())
	output := testRepo.Output().(*test.OutputLinear)
	oid := OidNow(output.Commits[0].Oid)

	// make sure this operation is valid
	_, err := repo.ObjectFromOid(oid)
	if err != nil {
		b.Errorf("could not read blob: %s", oid)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		repo.ObjectFromOid(oid)
	}
}
