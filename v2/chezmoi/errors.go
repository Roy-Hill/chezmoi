package chezmoi

import (
	"fmt"
	"os"
	"strings"
)

type duplicateTargetError struct {
	targetName  string
	sourcePaths []string
}

func (e *duplicateTargetError) Error() string {
	return fmt.Sprintf("%s: duplicate target (%s)", e.targetName, strings.Join(e.sourcePaths, ", "))
}

type unsupportedFileTypeError struct {
	path string
	mode os.FileMode
}

func (e *unsupportedFileTypeError) Error() string {
	return fmt.Sprintf("%s: unsupported file type %s", e.path, modeTypeName(e.mode))
}
