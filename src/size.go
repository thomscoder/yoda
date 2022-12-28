package src

import (
	"os"
	"path/filepath"
)

// dirSize returns the size of the given directory in bytes.
func dirSize(dir string) (int, error) {
	var size int
	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += int(info.Size())
		}
		return nil
	})
	return size, err
}
