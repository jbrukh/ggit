package ggit

// ================================================================= //
// FILE MODE
// ================================================================= //

type FileMode uint16

const (
    // TODO: there are more modes
    ModeDeleted    FileMode = 0000000
    ModeFile       FileMode = 0100644
    ModeExecutable FileMode = 0100755
    ModeTree       FileMode = 0040000
    ModeLink       FileMode = 0120000
)

func assertFileMode(u uint16) (FileMode, bool) {
    m := FileMode(u)
    switch m {
    case ModeDeleted,
        ModeExecutable,
        ModeFile,
        ModeTree,
        ModeLink:
        return m, true
    }
    return 0, false
}
