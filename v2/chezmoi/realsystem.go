package chezmoi

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/bmatcuk/doublestar"
	"github.com/google/renameio"
	vfs "github.com/twpayne/go-vfs"
	"go.uber.org/multierr"
)

// An RealSystem is a System that writes to a filesystem and executes scripts.
type RealSystem struct {
	vfs.FS
	PersistentState
	devCache     map[string]uint // devCache maps directories to device numbers.
	tempDirCache map[uint]string // tempDir maps device numbers to renameio temporary directories.
}

// NewRealSystem returns a System that acts on fs.
func NewRealSystem(fs vfs.FS, persistentState PersistentState) *RealSystem {
	return &RealSystem{
		FS:              fs,
		PersistentState: persistentState,
		devCache:        make(map[string]uint),
		tempDirCache:    make(map[uint]string),
	}
}

// Glob implements System.Glob.
func (s *RealSystem) Glob(pattern string) ([]string, error) {
	return doublestar.GlobOS(s, pattern)
}

// IdempotentCmdOutput implements System.IdempotentCmdOutput.
func (s *RealSystem) IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error) {
	return cmd.Output()
}

// PathSeparator implements doublestar.OS.PathSeparator.
func (s *RealSystem) PathSeparator() rune {
	return pathSeparator
}

// RunScript implements System.RunScript.
func (s *RealSystem) RunScript(name string, data []byte) (err error) {
	// Write the temporary script file. Put the randomness at the front of the
	// filename to preserve any file extension for Windows scripts.
	f, err := ioutil.TempFile("", "*."+path.Base(name))
	if err != nil {
		return
	}
	defer func() {
		err = multierr.Append(err, os.RemoveAll(f.Name()))
	}()

	// Make the script private before writing it in case it contains any
	// secrets.
	if err = f.Chmod(0o700); err != nil {
		return
	}
	_, err = f.Write(data)
	err = multierr.Append(err, f.Close())
	if err != nil {
		return
	}

	// Run the temporary script file.
	//nolint:gosec
	c := exec.Command(f.Name())
	// c.Dir = path.Join(applyOptions.DestDir, filepath.Dir(s.targetName)) // FIXME
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	return
}

// WriteSymlink implements System.WriteSymlink.
func (s *RealSystem) WriteSymlink(oldname, newname string) error {
	// Special case: if writing to the real filesystem, use
	// github.com/google/renameio.
	if s.FS == vfs.OSFS {
		return renameio.Symlink(oldname, newname)
	}
	if err := s.FS.RemoveAll(newname); err != nil && !os.IsNotExist(err) {
		return err
	}
	return s.FS.Symlink(oldname, newname)
}

// WriteFile is like ioutil.WriteFile but always sets perm before writing data.
// ioutil.WriteFile only sets the permissions when creating a new file. We need
// to ensure permissions, so we use our own implementation.
func WriteFile(fs vfs.FS, filename string, data []byte, perm os.FileMode) (err error) {
	// Create a new file, or truncate any existing one.
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return
	}
	defer func() {
		err = multierr.Append(err, f.Close())
	}()

	// Set permissions after truncation but before writing any data, in case the
	// file contained private data before, but before writing the new contents,
	// in case the contents contain private data after.
	if err = f.Chmod(perm); err != nil {
		return
	}

	_, err = f.Write(data)
	return err
}
