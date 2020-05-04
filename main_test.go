//+build !windows

// FIXME enable integration tests on Windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

//nolint:interfacer
func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"chezmoi": testRun,
	}))
}

func TestChezmoi(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("testdata", "scripts"),
		Setup: func(e *testscript.Env) error {
			var (
				homeDir        = filepath.Join(e.WorkDir, "home", "user")
				usrLocalBinDir = filepath.Join(e.WorkDir, "usr", "local", "bin")
			)

			for path, contents := range map[string][]string{
				// .gitconfig is populated with a user and email to avoid
				// warnings from git.
				filepath.Join(homeDir, ".gitconfig"): {
					"[user]",
					"    name = John Smith",
					"    email = john@home.org",
				},
				// editor a non-interactive script that appends "# edited\n" to
				// the end of each file.
				filepath.Join(usrLocalBinDir, "editor"): {
					"#!/bin/sh",
					"",
					"for filename in $*; do",
					"    echo \"# edited\" >> $filename",
					"done",
				},
			} {
				if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path, []byte(strings.Join(contents, "\n")), 0o777); err != nil {
					return err
				}
			}

			e.Setenv("EDITOR", "editor")
			e.Setenv("HOME", homeDir)
			e.Setenv("PATH", usrLocalBinDir+string(os.PathListSeparator)+e.Getenv("PATH"))

			return nil
		},
	})
}

func testRun() int {
	if err := run(); err != nil {
		if s := err.Error(); s != "" {
			fmt.Printf("chezmoi: %s\n", s)
		}
		return 1
	}
	return 0
}
