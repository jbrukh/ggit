package ggit

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

func assertFileMode(u uint16) (FileMode, bool) {
	m := FileMode(u)
	switch m {
	case ModeNew,
		ModeTree,
		ModeBlob,
		ModeBlobExec,
		ModeLink,
		ModeCommit:
		return m, true
	}
	return 0, false
}

// ================================================================= //
// OBJECT PARSER FUNCTIONS FOR FILE MODE
// ================================================================= //

func (p *objectParser) ParseFileMode(delim byte) (mode FileMode) {
	var ok bool
	if mode, ok = assertFileMode(uint16(p.ParseInt(delim, 8, 32))); !ok {
		panicErrf("expected: filemode")
	}
	return
}
