package chezmoi

import (
	"archive/tar"
	"io"
	"os"
	"os/user"
	"strconv"
	"time"
)

// A TARSystem is a System that writes to a TAR archive.
type TARSystem struct {
	nullSystem
	w              *tar.Writer
	headerTemplate tar.Header
}

// NewTARSystem returns a new TARSystem that writes a TAR file to w.
func NewTARSystem(w io.Writer, headerTemplate tar.Header) *TARSystem {
	return &TARSystem{
		w:              tar.NewWriter(w),
		headerTemplate: headerTemplate,
	}
}

// Chmod implements System.Chmod.
func (s *TARSystem) Chmod(name string, mode os.FileMode) error {
	return os.ErrPermission
}

// Close closes m.
func (s *TARSystem) Close() error {
	return s.w.Close()
}

// Mkdir implements System.Mkdir.
func (s *TARSystem) Mkdir(name string, perm os.FileMode) error {
	header := s.headerTemplate
	header.Typeflag = tar.TypeDir
	header.Name = name + "/"
	header.Mode = int64(perm)
	return s.w.WriteHeader(&header)
}

// RemoveAll implements System.RemoveAll.
func (s *TARSystem) RemoveAll(name string) error {
	return os.ErrPermission
}

// Rename implements System.Rename.
func (s *TARSystem) Rename(oldpath, newpath string) error {
	return os.ErrPermission
}

// RunScript implements System.RunScript.
func (s *TARSystem) RunScript(name string, data []byte) error {
	return s.WriteFile(name, data, 0o700)
}

// WriteFile implements System.WriteFile.
func (s *TARSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	header := s.headerTemplate
	header.Typeflag = tar.TypeReg
	header.Name = filename
	header.Size = int64(len(data))
	header.Mode = int64(perm)
	if err := s.w.WriteHeader(&header); err != nil {
		return err
	}
	_, err := s.w.Write(data)
	return err
}

// WriteSymlink implements System.WriteSymlink.
func (s *TARSystem) WriteSymlink(oldname, newname string) error {
	header := s.headerTemplate
	header.Typeflag = tar.TypeSymlink
	header.Name = newname
	header.Linkname = oldname
	return s.w.WriteHeader(&header)
}

// TARHeaderTemplate returns a tar.Header template populated with the current
// user and time.
func TARHeaderTemplate() tar.Header {
	// Attempt to lookup the current user. Ignore errors because the default
	// zero values are reasonable.
	var (
		uid   int
		gid   int
		Uname string
		Gname string
	)
	if currentUser, err := user.Current(); err == nil {
		uid, _ = strconv.Atoi(currentUser.Uid)
		gid, _ = strconv.Atoi(currentUser.Gid)
		Uname = currentUser.Username
		if group, err := user.LookupGroupId(currentUser.Gid); err == nil {
			Gname = group.Name
		}
	}

	now := time.Now()
	return tar.Header{
		Uid:        uid,
		Gid:        gid,
		Uname:      Uname,
		Gname:      Gname,
		ModTime:    now,
		AccessTime: now,
		ChangeTime: now,
	}
}
