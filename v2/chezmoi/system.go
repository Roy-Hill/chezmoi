package chezmoi

import (
	"os"
	"os/exec"
)

// A SystemReader reads from a filesystem and executes idempotent commands.
type SystemReader interface {
	PersistentStateReader
	Glob(pattern string) ([]string, error)
	IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error)
	Lstat(filename string) (os.FileInfo, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Readlink(name string) (string, error)
	Stat(name string) (os.FileInfo, error)
}

// A System reads from and writes to a filesystem, executes idempotent commands,
// runs scripts, and persists state.
type System interface {
	PersistentStateWriter
	SystemReader
	Chmod(name string, mode os.FileMode) error
	Mkdir(name string, perm os.FileMode) error
	RemoveAll(name string) error
	Rename(oldpath, newpath string) error
	RunScript(name string, data []byte) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
	WriteSymlink(oldname, newname string) error
}
