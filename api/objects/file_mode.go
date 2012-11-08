//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
file_mode.go implements the git-supported file modes used mainly in trees.
*/
package objects

// ================================================================= //
// FILE MODE
// ================================================================= //

type FileMode uint16

const (
	ModeNew      FileMode = 0000000
	ModeTree     FileMode = 0040000
	ModeBlob     FileMode = 0100644
	ModeBlobExec FileMode = 0100755
	ModeLink     FileMode = 0120000
	ModeCommit   FileMode = 0160000
)
