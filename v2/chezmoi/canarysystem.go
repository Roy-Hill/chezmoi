package chezmoi

import (
	"os"
	"os/exec"
)

// An CanarySystem wraps a System and records if any of its mutating
// methods are called.
type CanarySystem struct {
	s       System
	mutated bool
}

// NewCanarySystem returns a new CanarySystem.
func NewCanarySystem(s System) *CanarySystem {
	return &CanarySystem{
		s:       s,
		mutated: false,
	}
}

// Chmod implements System.Chmod.
func (s *CanarySystem) Chmod(name string, mode os.FileMode) error {
	s.mutated = true
	return s.s.Chmod(name, mode)
}

// Delete implements System.Delete.
func (s *CanarySystem) Delete(bucket, key []byte) error {
	s.mutated = true
	return s.s.Delete(bucket, key)
}

// Get implements System.Get.
func (s *CanarySystem) Get(bucket, key []byte) ([]byte, error) {
	return s.s.Get(bucket, key)
}

// Glob implements System.Glob.
func (s *CanarySystem) Glob(pattern string) ([]string, error) {
	return s.s.Glob(pattern)
}

// IdempotentCmdOutput implements System.IdempotentCmdOutput.
func (s *CanarySystem) IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error) {
	return s.s.IdempotentCmdOutput(cmd)
}

// Mkdir implements System.Mkdir.
func (s *CanarySystem) Mkdir(name string, perm os.FileMode) error {
	s.mutated = true
	return s.s.Mkdir(name, perm)
}

// Lstat implements System.Lstat.
func (s *CanarySystem) Lstat(path string) (os.FileInfo, error) {
	return s.s.Lstat(path)
}

// Mutated returns true if any of its mutating methods have been called.
func (s *CanarySystem) Mutated() bool {
	return s.mutated
}

// ReadDir implements System.ReadDir.
func (s *CanarySystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return s.s.ReadDir(dirname)
}

// ReadFile implements System.ReadFile.
func (s *CanarySystem) ReadFile(filename string) ([]byte, error) {
	return s.s.ReadFile(filename)
}

// Readlink implements System.Readlink.
func (s *CanarySystem) Readlink(name string) (string, error) {
	return s.s.Readlink(name)
}

// RemoveAll implements System.RemoveAll.
func (s *CanarySystem) RemoveAll(name string) error {
	s.mutated = true
	return s.s.RemoveAll(name)
}

// Rename implements System.Rename.
func (s *CanarySystem) Rename(oldpath, newpath string) error {
	s.mutated = true
	return s.s.Rename(oldpath, newpath)
}

// RunScript implements System.RunScript.
func (s *CanarySystem) RunScript(name string, data []byte) error {
	s.mutated = true
	return s.s.RunScript(name, data)
}

// Set implements System.Set.
func (s *CanarySystem) Set(bucket, key, value []byte) error {
	s.mutated = true
	return s.s.Set(bucket, key, value)
}

// Stat implements System.Stat.
func (s *CanarySystem) Stat(path string) (os.FileInfo, error) {
	return s.s.Stat(path)
}

// WriteFile implements System.WriteFile.
func (s *CanarySystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	s.mutated = true
	return s.s.WriteFile(name, data, perm)
}

// WriteSymlink implements System.WriteSymlink.
func (s *CanarySystem) WriteSymlink(oldname, newname string) error {
	s.mutated = true
	return s.s.WriteSymlink(oldname, newname)
}
