package util

import (
	"os"
	"path/filepath"
	"strings"
)

// GetDirAllFilePaths gets all the file paths in the specified directory recursively.
func GetDirAllFilePaths(dirname string, prefix string, suffix string) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := filepath.Join(dirname, info.Name())
		if info.IsDir() {
			tmp, err := GetDirAllFilePaths(path, prefix, suffix)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}

		var equal = true
		if suffix != "" && !strings.HasSuffix(path, suffix) {
			equal = false
		}
		if equal && prefix != "" && !strings.HasPrefix(path, prefix) {
			equal = false
		}

		if equal {
			paths = append(paths, path)
		}
	}
	return paths, nil
}
