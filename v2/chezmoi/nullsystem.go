package chezmoi

import (
	"os"
	"os/exec"
)

// An nullSystem represents an null System.
type nullSystem struct {
	nullPersistentState
}

func (nullSystem) Chmod(name string, mode os.FileMode) error                      { return nil }
func (nullSystem) Glob(pattern string) ([]string, error)                          { return nil, nil }
func (nullSystem) IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error)              { return cmd.Output() }
func (nullSystem) Lstat(name string) (os.FileInfo, error)                         { return nil, os.ErrNotExist }
func (nullSystem) Mkdir(dirname string, perm os.FileMode) error                   { return nil }
func (nullSystem) ReadDir(dirname string) ([]os.FileInfo, error)                  { return nil, os.ErrNotExist }
func (nullSystem) ReadFile(filename string) ([]byte, error)                       { return nil, os.ErrNotExist }
func (nullSystem) Readlink(name string) (string, error)                           { return "", os.ErrNotExist }
func (nullSystem) RemoveAll(name string) error                                    { return nil }
func (nullSystem) Rename(oldpath, newpath string) error                           { return nil }
func (nullSystem) RunScript(name string, data []byte) error                       { return nil }
func (nullSystem) Stat(name string) (os.FileInfo, error)                          { return nil, os.ErrNotExist }
func (nullSystem) WriteFile(filename string, data []byte, perm os.FileMode) error { return nil }
func (nullSystem) WriteSymlink(oldname, newname string) error                     { return nil }
