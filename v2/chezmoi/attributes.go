package chezmoi

import (
	"strings"
)

// A SourceFileType is a source file type.
type SourceFileType int

// Source file types.
const (
	SourceFileTypeFile SourceFileType = iota
	SourceFileTypeScript
	SourceFileTypeSymlink
)

// DirAttributes holds attributes parsed from a source directory name.
type DirAttributes struct {
	Name    string
	Exact   bool
	Private bool
}

// A FileAttributes holds attributes parsed from a source file name.
type FileAttributes struct {
	Name       string
	Type       SourceFileType
	Empty      bool
	Encrypted  bool
	Executable bool
	Once       bool
	Private    bool
	Template   bool
}

// ParseDirAttributes parses a single directory name.
func ParseDirAttributes(sourceName string) DirAttributes {
	var (
		name    = sourceName
		exact   = false
		private = false
	)
	if strings.HasPrefix(name, exactPrefix) {
		name = strings.TrimPrefix(name, exactPrefix)
		exact = true
	}
	if strings.HasPrefix(name, privatePrefix) {
		name = strings.TrimPrefix(name, privatePrefix)
		private = true
	}
	if strings.HasPrefix(name, dotPrefix) {
		name = "." + strings.TrimPrefix(name, dotPrefix)
	}
	return DirAttributes{
		Name:    name,
		Exact:   exact,
		Private: private,
	}
}

// SourceName returns da's source name.
func (da DirAttributes) SourceName() string {
	sourceName := ""
	if da.Exact {
		sourceName += exactPrefix
	}
	if da.Private {
		sourceName += privatePrefix
	}
	if strings.HasPrefix(da.Name, ".") {
		sourceName += dotPrefix + strings.TrimPrefix(da.Name, ".")
	} else {
		sourceName += da.Name
	}
	return sourceName
}

// ParseFileAttributes parses a source file name.
func ParseFileAttributes(sourceName string) FileAttributes {
	var (
		name       = sourceName
		typ        = SourceFileTypeFile
		empty      = false
		encrypted  = false
		executable = false
		once       = false
		private    = false
		template   = false
	)
	switch {
	case strings.HasPrefix(name, runPrefix):
		name = strings.TrimPrefix(name, runPrefix)
		typ = SourceFileTypeScript
		if strings.HasPrefix(name, oncePrefix) {
			name = strings.TrimPrefix(name, oncePrefix)
			once = true
		}
	case strings.HasPrefix(name, symlinkPrefix):
		name = strings.TrimPrefix(name, symlinkPrefix)
		typ = SourceFileTypeSymlink
		if strings.HasPrefix(name, dotPrefix) {
			name = "." + strings.TrimPrefix(name, dotPrefix)
		}
	default:
		if strings.HasPrefix(name, encryptedPrefix) {
			name = strings.TrimPrefix(name, encryptedPrefix)
			encrypted = true
		}
		if strings.HasPrefix(name, privatePrefix) {
			name = strings.TrimPrefix(name, privatePrefix)
			private = true
		}
		if strings.HasPrefix(name, emptyPrefix) {
			name = strings.TrimPrefix(name, emptyPrefix)
			empty = true
		}
		if strings.HasPrefix(name, executablePrefix) {
			name = strings.TrimPrefix(name, executablePrefix)
			executable = true
		}
		if strings.HasPrefix(name, dotPrefix) {
			name = "." + strings.TrimPrefix(name, dotPrefix)
		}
	}
	if strings.HasSuffix(name, templateSuffix) {
		name = strings.TrimSuffix(name, templateSuffix)
		template = true
	}
	return FileAttributes{
		Name:       name,
		Type:       typ,
		Empty:      empty,
		Encrypted:  encrypted,
		Executable: executable,
		Once:       once,
		Private:    private,
		Template:   template,
	}
}

// SourceName returns fa's source name.
func (fa FileAttributes) SourceName() string {
	sourceName := ""
	switch fa.Type {
	case SourceFileTypeFile:
		if fa.Encrypted {
			sourceName += encryptedPrefix
		}
		if fa.Private {
			sourceName += privatePrefix
		}
		if fa.Empty {
			sourceName += emptyPrefix
		}
		if fa.Executable {
			sourceName += executablePrefix
		}
	case SourceFileTypeScript:
		sourceName = runPrefix
		if fa.Once {
			sourceName += oncePrefix
		}
	case SourceFileTypeSymlink:
		sourceName = symlinkPrefix
	}
	if strings.HasPrefix(fa.Name, ".") {
		sourceName += dotPrefix + strings.TrimPrefix(fa.Name, ".")
	} else {
		sourceName += fa.Name
	}
	if fa.Template {
		sourceName += templateSuffix
	}
	return sourceName
}
