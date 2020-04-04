package chezmoi

import (
	"testing"

	"github.com/muesli/combinator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirAttributes(t *testing.T) {
	testData := struct {
		Name    []string
		Exact   []bool
		Private []bool
	}{
		Name:    []string{"dir", ".dir"},
		Exact:   []bool{false, true},
		Private: []bool{false, true},
	}
	var dirAttributes []DirAttributes
	require.NoError(t, combinator.Generate(&dirAttributes, testData))
	for _, da := range dirAttributes {
		actualSourceName := da.SourceName()
		actualDA := ParseDirAttributes(actualSourceName)
		assert.Equal(t, da, actualDA)
		assert.Equal(t, actualSourceName, actualDA.SourceName())
	}
}

func TestFileAttributes(t *testing.T) {
	var fileAttributes []FileAttributes
	require.NoError(t, combinator.Generate(&fileAttributes, struct {
		Type       []SourceFileType
		Name       []string
		Empty      []bool
		Encrypted  []bool
		Executable []bool
		Private    []bool
		Template   []bool
	}{
		Type:       []SourceFileType{SourceFileTypeFile},
		Name:       []string{"name", ".name"},
		Empty:      []bool{false, true},
		Encrypted:  []bool{false, true},
		Executable: []bool{false, true},
		Private:    []bool{false, true},
		Template:   []bool{false, true},
	}))
	require.NoError(t, combinator.Generate(&fileAttributes, struct {
		Type []SourceFileType
		Name []string
		Once []bool
	}{
		Type: []SourceFileType{SourceFileTypeScript},
		Name: []string{"name"},
		Once: []bool{false, true},
	}))
	require.NoError(t, combinator.Generate(&fileAttributes, struct {
		Type []SourceFileType
		Name []string
		Once []bool
	}{
		Type: []SourceFileType{SourceFileTypeSymlink},
		Name: []string{".name", "name"},
	}))
	for _, fa := range fileAttributes {
		actualSourceName := fa.SourceName()
		actualFA := ParseFileAttributes(actualSourceName)
		assert.Equal(t, fa, actualFA)
		assert.Equal(t, actualSourceName, actualFA.SourceName())
	}
}
