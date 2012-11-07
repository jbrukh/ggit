//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
trees_test.go implements git-comparison tests for ggit tree parsing.
*/
package api

import (
	//"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
	"testing"
)

var testCasesTreeReprs = []string{
	`dc5c0113b3d246da13b9d4b54861d2894da45e5a
348
100644 blob 022f0ad6583f99c3c01d0245db452e58bbb8fc8a	.gitignore
100644 blob 4760db634060e0f5c76eddb18031d53bff882a23	README
040000 tree 888fd8c9a6612c01d8f7c4e71fef114c3d3386d5	api
040000 tree 7ba590c4aa02ba2268f732f70b6143cd60457782	builtin
040000 tree bd52eecc257922f4c9008d080ad7dd5f9ac3444e	dotfiles
040000 tree 9ee2de084ec91660183740dacd8988c8170537e3	ggcase
100644 blob 82a1896a6d17e56220681fb66f94ac8a7e1d9bdc	ggit.go
100644 blob 3b20433fdf8544b2cff26989c5fc63327137e809	templates.go
040000 tree d42d8e7cf260b951e17687404031b1975ffa0b38	test
100644 blob 0c01aeb4856324d8487a252e32e04961bf603fad	version.go
`,

	`888fd8c9a6612c01d8f7c4e71fef114c3d3386d5
1240
100644 blob df0d1663548b5eaacc40819dc2cfc76c541e5fe3	blob.go
100644 blob b255d83c958a39352ec35eafc4355e5f35061e38	blob_test.go
100644 blob 6f7785c5a34aa603c2454c3e59f06244bd436319	commit.go
100644 blob 44bc403cc1c38122acb70fc59c37f353fee5b593	commit_test.go
100644 blob 94e4410342bb5c8a4aa2ecafce4c18281507a1c6	file_mode.go
100644 blob b9121f67f45bb44d40c41cb082e2a4247ce38fbc	file_mode_test.go
100644 blob 668cc6c8c2c900ca7e3374e43fc92b9986f1bba7	filter.go
100644 blob 6156f61d21b2fb2d17b486dd6564743d87eb7c21	filter_test.go
100644 blob e5d34a914a79b42eb0d22cd6e5337b49c027663e	format.go
100644 blob 644a1cb29278d81e5845afd2c8acb85a6e3bfe3f	format_test.go
100644 blob 6e5a89f2885ff4bfa12de9b56e5ce6308affddd5	index.go
100644 blob d6f665a9c18203dfb76640d897c89fbdb6946421	object.go
100644 blob 60d9a93f1696eb7b5dde56d0c69e805a7caee470	object_id.go
100644 blob 3c5843495b0ad6f7bcbed80c65fdf7876620af5b	object_id_test.go
100644 blob 3d9d4e964f3f4080f5dc4430ffeb2cc160967a8b	object_parser.go
100644 blob 0855c1b82077956a0d183c73c49d5d6bfd9b81e4	object_parser_test.go
100644 blob f6c0c08eda94ff37d10e8b10f209b00011ed99ef	object_type.go
100644 blob 6cad3f0da59bd4e1d5db77a30fa6a8798442007e	parser.go
100644 blob 16f62e0d7c5729d03d7027478e4b253e746e11bd	parser_test.go
100644 blob 51e81f899d1b4aa6cbc428529be72c6a231b5124	refs.go
100644 blob ec2087ab40b7cf9ddb97c585130399a8462cebda	refs_test.go
100644 blob 668158c18d3f7b4ab27b2d96a8d0cecdc1abb9c0	repository.go
100644 blob 778f64ec17cd4fd767e18d43231361d3aff70366	repository_test.go
100644 blob e75c23e6c6edfd17e9c7187014658c71d8e8761a	rev_parser.go
100644 blob 196267efc6c876663d68463fc9c70c34ead836a2	rev_parser_test.go
100644 blob 3055a0562b6324c6e01cd7c137628a9ea9efe141	tag.go
100644 blob 4277178f10be15a7dbf7e3cc2fb11ef730f72f21	tag_test.go
100644 blob 15bf87e3a4f2c55d59fd8532abc8d9aa25bacea4	tree.go
100644 blob d1e07986de64497fdbc4ba198c328b2852612727	util.go
100644 blob 778f64ec17cd4fd767e18d43231361d3aff70366	util_test.go
100644 blob 84dbccd834b954b24b5b0b4c56e1f1e4c1b3283a	who_when.go
`,
}

func Test_parseAndCompareTest(t *testing.T) {
	repo := Open("./..")
	for _, treeRepr := range testCasesTreeReprs {
		parseAndCompareTree(t, repo, treeRepr)
	}
}

// parseAndCompareTree will take a string tree representation and compare
// it with that tree from the provided repository. A tree representation
// has the following format: 
//
//     <oid_of_this_tree><LF>
//     <size><LF>
//     <output_of_cat_file_p>
//
// Example:
//
// `dc5c0113b3d246da13b9d4b54861d2894da45e5a
// 348
// 100644 blob 022f0ad6583f99c3c01d0245db452e58bbb8fc8a	.gitignore
// 100644 blob 4760db634060e0f5c76eddb18031d53bff882a23	README
// 040000 tree 888fd8c9a6612c01d8f7c4e71fef114c3d3386d5	api
// 040000 tree 7ba590c4aa02ba2268f732f70b6143cd60457782	builtin
// 040000 tree bd52eecc257922f4c9008d080ad7dd5f9ac3444e	dotfiles
// 040000 tree 9ee2de084ec91660183740dacd8988c8170537e3	ggcase
// 100644 blob 82a1896a6d17e56220681fb66f94ac8a7e1d9bdc	ggit.go
// 100644 blob 3b20433fdf8544b2cff26989c5fc63327137e809	templates.go
// 040000 tree d42d8e7cf260b951e17687404031b1975ffa0b38	test
// 100644 blob 0c01aeb4856324d8487a252e32e04961bf603fad	version.go
// `
func parseAndCompareTree(t *testing.T, repo Repository, treeRepr string) {
	p := objectParserForString(treeRepr)

	oid := p.ParseOid()
	p.ConsumeByte(LF)
	size := p.ParseInt(LF, 10, 32)

	o, err := repo.ObjectFromOid(oid)
	util.AssertNoErr(t, err)

	hdr := o.Header()
	util.Assert(t, hdr.Type() == objects.ObjectTree)
	util.Assert(t, hdr.Size() == size)

	tree := o.(*Tree)
	entries := tree.Entries()

	for !p.EOF() {
		mode := p.ParseFileMode(SP)
		otype := objects.ObjectType(p.ConsumeStrings(objectTypes))
		p.ConsumeByte(SP)
		oidStr := p.ParseOid().String()
		p.ConsumeByte(TAB)
		name := p.ReadString(LF)
		util.Assert(t, len(entries) > 0)
		entry := entries[0]
		util.Assert(t, mode == entry.mode)
		util.Assertf(t, name == entry.name, "expecting: `%s` got: `%s", name, entry.name)
		util.Assert(t, otype == entry.otype)
		util.Assertf(t, oidStr == entry.oid.String(), "expecting: `%s` got: `%s", oidStr, entry.oid)
		entries = entries[1:]
	}
	util.Assert(t, len(entries) == 0)
}
