package chezmoi

import (
	"log"
	"os"
	"os/exec"
	"time"
)

// A DebugSystem wraps a System and logs all of the actions it executes.
type DebugSystem struct {
	s System
}

// NewDebugSystem returns a new DebugSystem.
func NewDebugSystem(s System) *DebugSystem {
	return &DebugSystem{
		s: s,
	}
}

// Chmod implements System.Chmod.
func (s *DebugSystem) Chmod(name string, mode os.FileMode) error {
	return Debugf("Chmod(%q, 0o%o)", []interface{}{name, mode}, func() error {
		return s.s.Chmod(name, mode)
	})
}

// Delete implements System.Delete.
func (s *DebugSystem) Delete(bucket, key []byte) error {
	return Debugf("Delete(%q, %q)", []interface{}{string(bucket), string(key)}, func() error {
		return s.s.Delete(bucket, key)
	})
}

// Get implements System.Get.
func (s *DebugSystem) Get(bucket, key []byte) ([]byte, error) {
	var value []byte
	err := Debugf("Get(%q, %q)", []interface{}{string(bucket), string(key)}, func() error {
		var err error
		value, err = s.s.Get(bucket, key)
		return err
	})
	return value, err
}

// Glob implements System.Glob.
func (s *DebugSystem) Glob(name string) ([]string, error) {
	var matches []string
	err := Debugf("Glob(%q)", []interface{}{name}, func() error {
		var err error
		matches, err = s.s.Glob(name)
		return err
	})
	return matches, err
}

// IdempotentCmdOutput implements System.IdempotentCmdOutput.
func (s *DebugSystem) IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error) {
	var output []byte
	cmdStr := ShellQuoteArgs(append([]string{cmd.Path}, cmd.Args[1:]...))
	err := Debugf("IdempotentCmdOutput(%q)", []interface{}{cmdStr}, func() error {
		var err error
		output, err = s.s.IdempotentCmdOutput(cmd)
		return err
	})
	return output, err
}

// Lstat implements System.Lstat.
func (s *DebugSystem) Lstat(name string) (os.FileInfo, error) {
	var info os.FileInfo
	err := Debugf("Lstat(%q)", []interface{}{name}, func() error {
		var err error
		info, err = s.s.Lstat(name)
		return err
	})
	return info, err
}

// Mkdir implements System.Mkdir.
func (s *DebugSystem) Mkdir(name string, perm os.FileMode) error {
	return Debugf("Mkdir(%q, 0o%o)", []interface{}{name, perm}, func() error {
		return s.s.Mkdir(name, perm)
	})
}

// ReadDir implements System.ReadDir.
func (s *DebugSystem) ReadDir(name string) ([]os.FileInfo, error) {
	var infos []os.FileInfo
	err := Debugf("ReadDir(%q)", []interface{}{name}, func() error {
		var err error
		infos, err = s.s.ReadDir(name)
		return err
	})
	return infos, err
}

// ReadFile implements System.ReadFile.
func (s *DebugSystem) ReadFile(filename string) ([]byte, error) {
	var data []byte
	err := Debugf("ReadFile(%q)", []interface{}{filename}, func() error {
		var err error
		data, err = s.s.ReadFile(filename)
		return err
	})
	return data, err
}

// Readlink implements System.Readlink.
func (s *DebugSystem) Readlink(name string) (string, error) {
	var linkname string
	err := Debugf("Readlink(%q)", []interface{}{name}, func() error {
		var err error
		linkname, err = s.s.Readlink(name)
		return err
	})
	return linkname, err
}

// RemoveAll implements System.RemoveAll.
func (s *DebugSystem) RemoveAll(name string) error {
	return Debugf("RemoveAll(%q)", []interface{}{name}, func() error {
		return s.s.RemoveAll(name)
	})
}

// Rename implements System.Rename.
func (s *DebugSystem) Rename(oldpath, newpath string) error {
	return Debugf("Rename(%q, %q)", []interface{}{oldpath, newpath}, func() error {
		return s.Rename(oldpath, newpath)
	})
}

// RunScript implements System.RunScript.
func (s *DebugSystem) RunScript(name string, data []byte) error {
	return Debugf("Run(%q, _)", []interface{}{name}, func() error {
		return s.s.RunScript(name, data)
	})
}

// Set implements System.Set.
func (s *DebugSystem) Set(bucket, key, value []byte) error {
	return Debugf("Set(%q, %q, %q)", []interface{}{string(bucket), string(key), string(value)}, func() error {
		return s.s.Set(bucket, key, value)
	})
}

// Stat implements System.Stat.
func (s *DebugSystem) Stat(name string) (os.FileInfo, error) {
	var info os.FileInfo
	err := Debugf("Stat(%q)", []interface{}{name}, func() error {
		var err error
		info, err = s.s.Stat(name)
		return err
	})
	return info, err
}

// WriteFile implements System.WriteFile.
func (s *DebugSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return Debugf("WriteFile(%q, _, 0%o, _)", []interface{}{name, perm}, func() error {
		return s.s.WriteFile(name, data, perm)
	})
}

// WriteSymlink implements System.WriteSymlink.
func (s *DebugSystem) WriteSymlink(oldname, newname string) error {
	return Debugf("WriteSymlink(%q, %q)", []interface{}{oldname, newname}, func() error {
		return s.s.WriteSymlink(oldname, newname)
	})
}

// Debugf logs debugging information about calling f.
func Debugf(format string, args []interface{}, f func() error) error {
	errChan := make(chan error)
	start := time.Now()
	go func(errChan chan<- error) {
		errChan <- f()
	}(errChan)
	select {
	case err := <-errChan:
		if err == nil {
			log.Printf(format+" (%s)", append(args, time.Since(start))...)
		} else {
			log.Printf(format+" == %v (%s)", append(args, err, time.Since(start))...)
		}
		return err
	case <-time.After(1 * time.Second):
		log.Printf(format, args...)
		err := <-errChan
		if err == nil {
			log.Printf(format+" (%s)", append(args, time.Since(start))...)
		} else {
			log.Printf(format+" == %v (%s)", append(args, err, time.Since(start))...)
		}
		return err
	}
}
