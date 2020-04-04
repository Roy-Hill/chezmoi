package chezmoi

import (
	"os"
	"testing"
	"text/template"

	"github.com/coreos/go-semver/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-vfs/vfst"
)

func TestSourceStateApplyAll(t *testing.T) {
	// FIXME script tests
	// FIXME script template tests
	for _, tc := range []struct {
		name               string
		root               interface{}
		sourceStateOptions []SourceStateOption
		tests              []interface{}
	}{
		{
			name: "empty",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": &vfst.Dir{Perm: 0o755},
				},
			},
		},
		{
			name: "dir",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": map[string]interface{}{
						"foo": &vfst.Dir{Perm: 0o755},
					},
				},
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestIsDir,
					vfst.TestModePerm(0o755),
				),
			},
		},
		{
			name: "dir_exact",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "",
					},
					".local/share/chezmoi": map[string]interface{}{
						"exact_foo": &vfst.Dir{Perm: 0o755},
					},
				},
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestIsDir,
					vfst.TestModePerm(0o755),
				),
				vfst.TestPath("/home/user/foo/bar",
					vfst.TestDoesNotExist,
				),
			},
		},
		{
			name: "file",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("bar"),
				),
			},
		},
		{
			name: "file_remove_empty",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					"foo": "",
					".local/share/chezmoi": map[string]interface{}{
						"foo": "",
					},
				},
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestDoesNotExist,
				),
			},
		},
		{
			name: "file_create_empty",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": map[string]interface{}{
						"empty_foo": "",
					},
				},
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString(""),
				),
			},
		},
		{
			name: "file_template",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": map[string]interface{}{
						"foo.tmpl": "email = {{ .email }}",
					},
				},
			},
			sourceStateOptions: []SourceStateOption{
				WithTemplateData(map[string]interface{}{
					"email": "john.smith@company.com",
				}),
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("email = john.smith@company.com"),
				),
			},
		},
		{
			name: "symlink",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": map[string]interface{}{
						"symlink_foo": "bar",
					},
				},
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestModeType(os.ModeSymlink),
					vfst.TestSymlinkTarget("bar"),
				),
			},
		},
		{
			name: "symlink_template",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					".local/share/chezmoi": map[string]interface{}{
						"symlink_foo.tmpl": "bar_{{ .os }}",
					},
				},
			},
			sourceStateOptions: []SourceStateOption{
				WithTemplateData(map[string]interface{}{
					"os": "linux",
				}),
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/foo",
					vfst.TestModeType(os.ModeSymlink),
					vfst.TestSymlinkTarget("bar_linux"),
				),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, cleanup, err := vfst.NewTestFS(tc.root)
			require.NoError(t, err)
			defer cleanup()

			s := newTestRealSystem(fs)
			sourceStateOptions := []SourceStateOption{
				WithSystem(s),
				WithSourcePath("/home/user/.local/share/chezmoi"),
			}
			sourceStateOptions = append(sourceStateOptions, tc.sourceStateOptions...)
			ss := NewSourceState(sourceStateOptions...)
			require.NoError(t, ss.Read())
			require.NoError(t, ss.Evaluate())
			require.NoError(t, ss.ApplyAll(s, vfst.DefaultUmask, "/home/user"))

			vfst.RunTests(t, fs, "", tc.tests...)
		})
	}
}

func TestSourceStateRead(t *testing.T) {
	for _, tc := range []struct {
		name                string
		root                interface{}
		sourceStateOptions  []SourceStateOption
		expectedError       string
		expectedSourceState *SourceState
	}{
		{
			name: "empty",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": &vfst.Dir{Perm: 0o755},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
			),
		},
		{
			name: "dir",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"foo": &vfst.Dir{Perm: 0o755},
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withEntries(map[string]SourceStateEntry{
					"foo": &SourceStateDir{
						path: "/home/user/.local/share/chezmoi/foo",
						attributes: DirAttributes{
							Name: "foo",
						},
						targetStateEntry: &TargetStateDir{
							perm: 0o755,
						},
					},
				}),
			),
		},
		{
			name: "file",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"foo": "bar",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withEntries(map[string]SourceStateEntry{
					"foo": &SourceStateFile{
						path: "/home/user/.local/share/chezmoi/foo",
						attributes: FileAttributes{
							Name: "foo",
							Type: SourceFileTypeFile,
						},
						lazyContents: newLazyContents([]byte("bar")),
						targetStateEntry: &TargetStateFile{
							perm:         0o644,
							lazyContents: newLazyContents([]byte("bar")),
						},
					},
				}),
			),
		},
		{
			name: "duplicate_target",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"foo":      "bar",
					"foo.tmpl": "bar",
				},
			},
			expectedError: "foo: duplicate target (/home/user/.local/share/chezmoi/foo, /home/user/.local/share/chezmoi/foo.tmpl)",
		},
		{
			name: "duplicate_target",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"foo":       "bar",
					"exact_foo": &vfst.Dir{Perm: 0o755},
				},
			},
			expectedError: "foo: duplicate target (/home/user/.local/share/chezmoi/exact_foo, /home/user/.local/share/chezmoi/foo)",
		},
		{
			name: "symlink",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"foo": &vfst.Symlink{Target: "bar"},
				},
			},
			expectedError: "/home/user/.local/share/chezmoi/foo: unsupported file type symlink",
		},
		{
			name: "script",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"run_foo": "bar",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withEntries(map[string]SourceStateEntry{
					"foo": &SourceStateFile{
						path: "/home/user/.local/share/chezmoi/run_foo",
						attributes: FileAttributes{
							Name: "foo",
							Type: SourceFileTypeScript,
						},
						lazyContents: newLazyContents([]byte("bar")),
						targetStateEntry: &TargetStateScript{
							name:         "foo",
							lazyContents: newLazyContents([]byte("bar")),
						},
					},
				}),
			),
		},
		{
			name: "symlink",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"symlink_foo": "bar",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withEntries(map[string]SourceStateEntry{
					"foo": &SourceStateFile{
						path: "/home/user/.local/share/chezmoi/symlink_foo",
						attributes: FileAttributes{
							Name: "foo",
							Type: SourceFileTypeSymlink,
						},
						lazyContents: newLazyContents([]byte("bar")),
						targetStateEntry: &TargetStateSymlink{
							lazyLinkname: newLazyLinkname("bar"),
						},
					},
				}),
			),
		},
		{
			name: "file_in_dir",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "baz",
					},
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withEntries(map[string]SourceStateEntry{
					"foo": &SourceStateDir{
						path: "/home/user/.local/share/chezmoi/foo",
						attributes: DirAttributes{
							Name: "foo",
						},
						targetStateEntry: &TargetStateDir{
							perm: 0o755,
						},
					},
					"foo/bar": &SourceStateFile{
						path: "/home/user/.local/share/chezmoi/foo/bar",
						attributes: FileAttributes{
							Name: "bar",
							Type: SourceFileTypeFile,
						},
						lazyContents: &lazyContents{
							contents: []byte("baz"),
						},
						targetStateEntry: &TargetStateFile{
							perm: 0o644,
							lazyContents: &lazyContents{
								contents: []byte("baz"),
							},
						},
					},
				}),
			),
		},
		{
			name: "chezmoiignore",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".chezmoiignore": "README.md\n",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withIgnore(
					NewPatternSet(
						withAdd(t, "README.md", true),
					),
				),
			),
		},
		{
			name: "chezmoiignore_ignore_file",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".chezmoiignore": "README.md\n",
					"README.md":      "",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withIgnore(
					NewPatternSet(
						withAdd(t, "README.md", true),
					),
				),
			),
		},
		{
			name: "chezmoiremove",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".chezmoiremove": "!*.txt\n",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withRemove(
					NewPatternSet(
						withAdd(t, "*.txt", false),
					),
				),
			),
		},
		{
			name: "chezmoitemplates",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".chezmoitemplates": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withTemplates(
					map[string]*template.Template{
						"foo": template.Must(template.New("foo").Parse("bar")),
					},
				),
			),
		},
		{
			name: "chezmoiversion",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".chezmoiversion": "1.2.3\n",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withMinVersion(
					&semver.Version{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
				),
			),
		},
		{
			name: "chezmoiversion_multiple",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".chezmoiversion": "1.2.3\n",
					"foo": map[string]interface{}{
						".chezmoiversion": "2.3.4\n",
					},
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
				withEntries(map[string]SourceStateEntry{
					"foo": &SourceStateDir{
						path: "/home/user/.local/share/chezmoi/foo",
						attributes: DirAttributes{
							Name: "foo",
						},
						targetStateEntry: &TargetStateDir{
							perm: 0o755,
						},
					},
				}),
				withMinVersion(
					&semver.Version{
						Major: 2,
						Minor: 3,
						Patch: 4,
					},
				),
			),
		},
		{
			name: "ignore_dir",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".ignore": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
			),
		},
		{
			name: "ignore_file",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": map[string]interface{}{
					".ignore": "",
				},
			},
			expectedSourceState: NewSourceState(
				WithSourcePath("/home/user/.local/share/chezmoi"),
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, cleanup, err := vfst.NewTestFS(tc.root)
			require.NoError(t, err)
			defer cleanup()

			sourceStateOptions := []SourceStateOption{
				WithSystem(newTestRealSystem(fs)),
				WithSourcePath("/home/user/.local/share/chezmoi"),
			}
			sourceStateOptions = append(sourceStateOptions, tc.sourceStateOptions...)
			ss := NewSourceState(sourceStateOptions...)
			err = ss.Read()
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
			require.NoError(t, ss.Evaluate())
			require.NoError(t, tc.expectedSourceState.Evaluate())
			ss.s = nil
			assert.Equal(t, tc.expectedSourceState, ss)
		})
	}
}

func TestSourceStateRemove(t *testing.T) {
	// FIXME test doublestar globs
	for _, tc := range []struct {
		name  string
		root  interface{}
		tests []vfst.Test
	}{
		{
			name: "empty",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi": &vfst.Dir{Perm: 0o755},
			},
		},
		{
			name: "dir",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					"dir":     &vfst.Dir{Perm: 0o755},
					"file":    "",
					"symlink": &vfst.Symlink{Target: "file"},
					".local/share/chezmoi": map[string]interface{}{
						".chezmoiremove": "dir\n",
					},
				},
			},
			tests: []vfst.Test{
				vfst.TestPath("/home/user/dir",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/file",
					vfst.TestModeIsRegular,
				),
				vfst.TestPath("/home/user/symlink",
					vfst.TestModeType(os.ModeSymlink),
				),
			},
		},
		{
			name: "file",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					"dir":     &vfst.Dir{Perm: 0o755},
					"file":    "",
					"symlink": &vfst.Symlink{Target: "file"},
					".local/share/chezmoi": map[string]interface{}{
						".chezmoiremove": "file\n",
					},
				},
			},
			tests: []vfst.Test{
				vfst.TestPath("/home/user/dir",
					vfst.TestIsDir,
				),
				vfst.TestPath("/home/user/file",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/symlink",
					vfst.TestModeType(os.ModeSymlink),
				),
			},
		},
		{
			name: "symlink",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					"dir":     &vfst.Dir{Perm: 0o755},
					"file":    "",
					"symlink": &vfst.Symlink{Target: "file"},
					".local/share/chezmoi": map[string]interface{}{
						".chezmoiremove": "symlink\n",
					},
				},
			},
			tests: []vfst.Test{
				vfst.TestPath("/home/user/dir",
					vfst.TestIsDir,
				),
				vfst.TestPath("/home/user/file",
					vfst.TestModeIsRegular,
				),
				vfst.TestPath("/home/user/symlink",
					vfst.TestDoesNotExist,
				),
			},
		},
		{
			name: "exclude_pattern",
			root: map[string]interface{}{
				"/home/user": map[string]interface{}{
					"foo": "",
					"bar": "",
					"baz": "",
					".local/share/chezmoi": map[string]interface{}{
						".chezmoiremove": "b*\n!*z\n",
					},
				},
			},
			tests: []vfst.Test{
				vfst.TestPath("/home/user/foo",
					vfst.TestModeIsRegular,
				),
				vfst.TestPath("/home/user/bar",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/baz",
					vfst.TestModeIsRegular,
				),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, cleanup, err := vfst.NewTestFS(tc.root)
			require.NoError(t, err)
			defer cleanup()

			ss := NewSourceState(
				WithSystem(newTestRealSystem(fs)),
				WithSourcePath("/home/user/.local/share/chezmoi"),
			)
			require.NoError(t, ss.Read())
			require.NoError(t, ss.Evaluate())

			require.NoError(t, ss.Remove(newTestRealSystem(fs), "/home/user"))

			vfst.RunTests(t, fs, "", tc.tests)
		})
	}
}

func withEntries(sourceEntries map[string]SourceStateEntry) SourceStateOption {
	return func(ss *SourceState) {
		ss.Entries = sourceEntries
	}
}

func withIgnore(ignore *PatternSet) SourceStateOption {
	return func(ss *SourceState) {
		ss.ignore = ignore
	}
}

func withMinVersion(minVersion *semver.Version) SourceStateOption {
	return func(ss *SourceState) {
		ss.minVersion = minVersion
	}
}

func withRemove(remove *PatternSet) SourceStateOption {
	return func(ss *SourceState) {
		ss.remove = remove
	}
}

func withTemplates(templates map[string]*template.Template) SourceStateOption {
	return func(ss *SourceState) {
		ss.templates = templates
	}
}
