package filesearch

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ErrDirectoryNotExist = errors.New("directory does not exist")
	ErrNoExtensions      = errors.New("no file extensions provided")
)

// SearchFiles recursive searches for files with the specified extensions in the given directory and returns a sorted list of file paths.
func SearchFiles(root string, extensions []string) ([]string, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, ErrDirectoryNotExist
	}

	if len(extensions) == 0 {
		return nil, ErrNoExtensions
	}

	var files []string

	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			for ext := range extMap {
				if strings.HasSuffix(info.Name(), ext) {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}
