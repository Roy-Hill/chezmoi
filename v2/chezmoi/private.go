// +build !windows

package chezmoi

import (
	"runtime"

	vfs "github.com/twpayne/go-vfs"
)

// IsPrivate returns whether path should be considered private.
func IsPrivate(fs vfs.Stater, path string, want bool) (bool, error) {
	// Private has no real equivalent on Windows, so always return what the
	// caller wants.
	if runtime.GOOS == "windows" {
		return want, nil
	}

	info, err := fs.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().Perm()&0o77 == 0, nil
}
