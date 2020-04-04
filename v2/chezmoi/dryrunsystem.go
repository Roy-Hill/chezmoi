package chezmoi

import (
	"os"
	"os/exec"
)

// DryRunSystem is an System that reads from, but does not write to, to
// a wrapped System.
type DryRunSystem struct {
	s System
}

// NewDryRunSystem returns a new DryRunSystem that wraps fs.
func NewDryRunSystem(s System) *DryRunSystem {
	return &DryRunSystem{
		s: s,
	}
}

// Chmod implements System.Chmod.
func (s *DryRunSystem) Chmod(name string, mode os.FileMode) error {
	return nil
}

// Delete implements System.Delete.
func (s *DryRunSystem) Delete(bucket, key []byte) error {
	return nil
}

// Get implements System.Get.
func (s *DryRunSystem) Get(bucket, key []byte) ([]byte, error) {
	return s.s.Get(bucket, key)
}

// Glob implements System.Glob.
func (s *DryRunSystem) Glob(pattern string) ([]string, error) {
	return s.s.Glob(pattern)
}

// IdempotentCmdOutput implements System.IdempotentCmdOutput.
func (s *DryRunSystem) IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error) {
	return s.s.IdempotentCmdOutput(cmd)
}

// Lstat implements System.Lstat.
func (s *DryRunSystem) Lstat(name string) (os.FileInfo, error) {
	return s.s.Stat(name)
}

// Mkdir implements System.Mkdir.
func (s *DryRunSystem) Mkdir(name string, perm os.FileMode) error {
	return nil
}

// ReadDir implements System.ReadDir.
func (s *DryRunSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return s.s.ReadDir(dirname)
}

// ReadFile implements System.ReadFile.
func (s *DryRunSystem) ReadFile(filename string) ([]byte, error) {
	return s.s.ReadFile(filename)
}

// Readlink implements System.Readlink.
func (s *DryRunSystem) Readlink(name string) (string, error) {
	return s.s.Readlink(name)
}

// RemoveAll implements System.RemoveAll.
func (s *DryRunSystem) RemoveAll(string) error {
	return nil
}

// Rename implements System.Rename.
func (s *DryRunSystem) Rename(oldpath, newpath string) error {
	return nil
}

// RunScript implements System.RunScript.
func (s *DryRunSystem) RunScript(name string, data []byte) error {
	return nil
}

// Set implements System.Set.
func (s *DryRunSystem) Set(bucket, key, value []byte) error {
	return nil
}

// Stat implements System.Stat.
func (s *DryRunSystem) Stat(name string) (os.FileInfo, error) {
	return s.s.Stat(name)
}

// WriteFile implements System.WriteFile.
func (s *DryRunSystem) WriteFile(string, []byte, os.FileMode) error {
	return nil
}

// WriteSymlink implements System.WriteSymlink.
func (s *DryRunSystem) WriteSymlink(string, string) error {
	return nil
}
