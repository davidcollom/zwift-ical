package site

import (
	"os"
	"strings"
)

// EnsureDir creates a directory if it does not exist
func EnsureDir(dir string) error {
	return os.MkdirAll(strings.ToLower(dir), 0755)
}
