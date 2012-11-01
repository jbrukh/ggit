//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_blobs.go implements a test repository.
*/
package test

// ================================================================= //
// TEST CASE: BUNCHES OF BLOBS
// ================================================================= //

type OutputBlob struct {
	Oid      string
	Contents string
}

type OutputBlobs struct {
	Blobs []*OutputBlob
	N     int
}

var testCasesBlobContents = []string{
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

var Blobs = NewRepoTestCase(
	"__blobs",
	func(testCase *RepoTestCase) (err error) {
		err = createRepo(testCase)
		if err != nil {
			return err
		}

		output := &OutputBlobs{
			Blobs: make([]*OutputBlob, 0),
			N:     len(testCasesBlobContents),
		}

		// hash the test objects
		for _, contents := range testCasesBlobContents {
			if oidStr, err := HashBlob(testCase.Repo(), contents); err != nil {
				return err
			} else {
				output.Blobs = append(output.Blobs, &OutputBlob{
					oidStr,
					contents,
				})
			}
		}
		testCase.output = output
		return
	},
)
