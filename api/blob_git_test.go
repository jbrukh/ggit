//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/util"
	"testing"
)

var blobContents = []string{
	`'Tis better to have loved and lost than never
	to     have
	lov3d
	@
	.*
	`,

	``,

	"This is a test.",

	"hahahahahaha",

	"49-230948fdskv93485ufdlskj3498",

	`На берегу пустынных волн
	Стоял он, дум великих полн,
	И вдаль глядел. Пред ним широко
	Река неслася; бедный чёлн
	По ней стремился одиноко.
	По мшистым, топким берегам
	Чернели избы здесь и там,
	Приют убогого чухонца;
	И лес, неведомый лучам
	В тумане спрятанного солнца,
	Кругом шумел.`,

	`யாமறிந்த மொழிகளிலே தமிழ்மொழி போல் இனிதாவது எங்கும் காணோம், 
	பாமரராய் விலங்குகளாய், உலகனைத்தும் இகழ்ச்சிசொலப் பான்மை கெட்டு, 
	நாமமது தமிழரெனக் கொண்டு இங்கு வாழ்ந்திடுதல் நன்றோ? சொல்லீர்!
	தேமதுரத் தமிழோசை உலகமெலாம் பரவும்வகை செய்தல் வேண்டும்.`,
}

func Test_readBlobs(t *testing.T) {
	repo := util.TempRepo("test_blobs")
	util.AssertCreateGitRepo(t, repo)
	defer util.AssertRemoveGitRepo(t, repo)

	// create a ggit repo
	ggrepo := Open(repo)

	// hash the test objects
	for _, contents := range blobContents {
		oidStr, err := util.HashBlob(repo, contents)
		util.AssertNoErr(t, err)
		oid := OidNow(oidStr)

		// read the blob
		o, err := ggrepo.ObjectFromOid(oid)
		util.AssertNoErr(t, err)
		util.Assert(t, o.Header().Type() == ObjectBlob)
		b := o.(*Blob)
		util.AssertEqualString(t, b.String(), contents)
		util.AssertEqualInt(t, b.Header().Size(), len(contents))
	}

}
