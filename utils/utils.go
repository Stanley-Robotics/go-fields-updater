package utils

import (
	"os"
)

// IsDir reports whether the named file is a directory.
func IsDir(name string) bool {
	info, err := os.Stat(name)
	return (err == nil) && info.IsDir()
}
