package chezmoi

import (
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// A GitDiffSystem wraps a SystemReader and logs all of the actions it
// would execute as a git diff.
type GitDiffSystem struct {
	sr             SystemReader
	prefix         string
	unifiedEncoder *diff.UnifiedEncoder
}

// NewGitDiffSystem returns a new GitDiffSystem.
func NewGitDiffSystem(unifiedEncoder *diff.UnifiedEncoder, sr SystemReader, prefix string) *GitDiffSystem {
	return &GitDiffSystem{
		sr:             sr,
		prefix:         prefix,
		unifiedEncoder: unifiedEncoder,
	}
}

// Chmod implements System.Chmod.
func (s *GitDiffSystem) Chmod(name string, mode os.FileMode) error {
	fromFileMode, info, err := s.getFileMode(name)
	if err != nil {
		return err
	}
	// Assume that we're only changing permissions.
	toFileMode, err := filemode.NewFromOSFileMode(info.Mode()&^os.ModePerm | mode)
	if err != nil {
		return err
	}
	path := s.trimPrefix(name)
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				from: &gitDiffFile{
					fileMode: fromFileMode,
					path:     path,
					hash:     plumbing.ZeroHash,
				},
				to: &gitDiffFile{
					fileMode: toFileMode,
					path:     path,
					hash:     plumbing.ZeroHash,
				},
			},
		},
	})
}

// Delete implements System.Delete.
func (s *GitDiffSystem) Delete(bucket, key []byte) error {
	return nil
}

// Get implements System.Get.
func (s *GitDiffSystem) Get(bucket, key []byte) ([]byte, error) {
	return s.sr.Get(bucket, key)
}

// Glob implements System.Glob.
func (s *GitDiffSystem) Glob(pattern string) ([]string, error) {
	return s.sr.Glob(pattern)
}

// Lstat implements System.Lstat.
func (s *GitDiffSystem) Lstat(name string) (os.FileInfo, error) {
	return s.sr.Lstat(name)
}

// IdempotentCmdOutput implements System.IdempotentCmdOutput.
func (s *GitDiffSystem) IdempotentCmdOutput(cmd *exec.Cmd) ([]byte, error) {
	return s.sr.IdempotentCmdOutput(cmd)
}

// Mkdir implements System.Mkdir.
func (s *GitDiffSystem) Mkdir(name string, perm os.FileMode) error {
	toFileMode, err := filemode.NewFromOSFileMode(os.ModeDir | perm)
	if err != nil {
		return err
	}
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				to: &gitDiffFile{
					fileMode: toFileMode,
					path:     s.trimPrefix(name),
					hash:     plumbing.ZeroHash,
				},
			},
		},
	})
}

// ReadDir implements System.ReadDir.
func (s *GitDiffSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return s.sr.ReadDir(dirname)
}

// ReadFile implements System.ReadFile.
func (s *GitDiffSystem) ReadFile(filename string) ([]byte, error) {
	return s.sr.ReadFile(filename)
}

// Readlink implements System.Readlink.
func (s *GitDiffSystem) Readlink(name string) (string, error) {
	return s.sr.Readlink(name)
}

// RemoveAll implements System.RemoveAll.
func (s *GitDiffSystem) RemoveAll(name string) error {
	fromFileMode, _, err := s.getFileMode(name)
	if err != nil {
		return err
	}
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				from: &gitDiffFile{
					fileMode: fromFileMode,
					path:     s.trimPrefix(name),
					hash:     plumbing.ZeroHash,
				},
			},
		},
	})
}

// RunScript implements System.RunScript.
func (s *GitDiffSystem) RunScript(name string, data []byte) error {
	isBinary := isBinary(data)
	var chunks []diff.Chunk
	if !isBinary {
		chunk := &gitDiffChunk{
			content:   string(data),
			operation: diff.Add,
		}
		chunks = append(chunks, chunk)
	}
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				isBinary: isBinary,
				to: &gitDiffFile{
					fileMode: filemode.Executable,
					path:     s.trimPrefix(name),
					hash:     plumbing.ComputeHash(plumbing.BlobObject, data),
				},
				chunks: chunks,
			},
		},
	})
}

// Stat implements System.Stat.
func (s *GitDiffSystem) Stat(name string) (os.FileInfo, error) {
	return s.sr.Stat(name)
}

// Rename implements System.Rename.
func (s *GitDiffSystem) Rename(oldpath, newpath string) error {
	fileMode, _, err := s.getFileMode(oldpath)
	if err != nil {
		return err
	}
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				from: &gitDiffFile{
					fileMode: fileMode,
					path:     s.trimPrefix(oldpath),
					hash:     plumbing.ZeroHash,
				},
				to: &gitDiffFile{
					fileMode: fileMode,
					path:     s.trimPrefix(newpath),
					hash:     plumbing.ZeroHash,
				},
			},
		},
	})
}

// Set implements System.Set.
func (s *GitDiffSystem) Set(bucket, key, value []byte) error {
	return nil
}

// WriteFile implements System.WriteFile.
func (s *GitDiffSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fromFileMode, _, err := s.getFileMode(filename)
	if err != nil {
		return err
	}
	fromData, err := s.sr.ReadFile(filename)
	if err != nil {
		return err
	}
	toFileMode, err := filemode.NewFromOSFileMode(perm)
	if err != nil {
		return err
	}
	path := s.trimPrefix(filename)
	isBinary := isBinary(fromData) || isBinary(data)
	var chunks []diff.Chunk
	if !isBinary {
		chunks = diffChunks(string(fromData), string(data))
	}
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				isBinary: isBinary,
				from: &gitDiffFile{
					fileMode: fromFileMode,
					path:     path,
					hash:     plumbing.ComputeHash(plumbing.BlobObject, fromData),
				},
				to: &gitDiffFile{
					fileMode: toFileMode,
					path:     path,
					hash:     plumbing.ComputeHash(plumbing.BlobObject, data),
				},
				chunks: chunks,
			},
		},
	})
}

// WriteSymlink implements System.WriteSymlink.
func (s *GitDiffSystem) WriteSymlink(oldname, newname string) error {
	return s.unifiedEncoder.Encode(&gitDiffPatch{
		filePatches: []diff.FilePatch{
			&gitDiffFilePatch{
				to: &gitDiffFile{
					fileMode: filemode.Symlink,
					path:     s.trimPrefix(newname),
					hash:     plumbing.ComputeHash(plumbing.BlobObject, []byte(oldname)),
				},
				chunks: []diff.Chunk{
					&gitDiffChunk{
						content:   oldname,
						operation: diff.Add,
					},
				},
			},
		},
	})
}

func (s *GitDiffSystem) getFileMode(name string) (filemode.FileMode, os.FileInfo, error) {
	info, err := s.sr.Stat(name)
	if err != nil {
		return filemode.Empty, nil, err
	}
	fileMode, err := filemode.NewFromOSFileMode(info.Mode())
	return fileMode, info, err
}

func (s *GitDiffSystem) trimPrefix(path string) string {
	return strings.TrimPrefix(path, s.prefix)
}

var gitDiffOperation = map[diffmatchpatch.Operation]diff.Operation{
	diffmatchpatch.DiffDelete: diff.Delete,
	diffmatchpatch.DiffEqual:  diff.Equal,
	diffmatchpatch.DiffInsert: diff.Add,
}

type gitDiffChunk struct {
	content   string
	operation diff.Operation
}

func (c *gitDiffChunk) Content() string      { return c.content }
func (c *gitDiffChunk) Type() diff.Operation { return c.operation }

type gitDiffFile struct {
	hash     plumbing.Hash
	fileMode filemode.FileMode
	path     string
}

func (f *gitDiffFile) Hash() plumbing.Hash     { return f.hash }
func (f *gitDiffFile) Mode() filemode.FileMode { return f.fileMode }
func (f *gitDiffFile) Path() string            { return f.path }

type gitDiffFilePatch struct {
	isBinary bool
	from, to diff.File
	chunks   []diff.Chunk
}

func (fp *gitDiffFilePatch) IsBinary() bool                { return fp.isBinary }
func (fp *gitDiffFilePatch) Files() (diff.File, diff.File) { return fp.from, fp.to }
func (fp *gitDiffFilePatch) Chunks() []diff.Chunk          { return fp.chunks }

type gitDiffPatch struct {
	filePatches []diff.FilePatch
	message     string
}

func (p *gitDiffPatch) FilePatches() []diff.FilePatch { return p.filePatches }
func (p *gitDiffPatch) Message() string               { return p.message }

func diffChunks(from, to string) []diff.Chunk {
	dmp := diffmatchpatch.New()
	dmp.DiffTimeout = time.Second
	fromRunes, toRunes, runesToLines := dmp.DiffLinesToRunes(from, to)
	diffs := dmp.DiffCharsToLines(dmp.DiffMainRunes(fromRunes, toRunes, false), runesToLines)
	chunks := make([]diff.Chunk, 0, len(diffs))
	for _, d := range diffs {
		chunk := &gitDiffChunk{
			content:   d.Text,
			operation: gitDiffOperation[d.Type],
		}
		chunks = append(chunks, chunk)
	}
	return chunks
}

func isBinary(data []byte) bool {
	return len(data) != 0 && !strings.HasPrefix(http.DetectContentType(data), "text/")
}
