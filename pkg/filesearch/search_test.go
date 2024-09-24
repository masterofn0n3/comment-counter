package filesearch

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDir() (string, error) {
	tempDir := os.TempDir()
	testDir := filepath.Join(tempDir, "testdir")

	os.RemoveAll(testDir)
	if err := os.Mkdir(testDir, 0755); err != nil {
		return "", err
	}

	subDirs := []string{
		"subdir1",
		"subdir2",
		"subdir2/subsubdir1",
		"subdir2/subsubdir2",
	}

	for _, dir := range subDirs {
		if err := os.MkdirAll(filepath.Join(testDir, dir), 0755); err != nil {
			return "", err
		}
	}

	files := []struct {
		path string
		name string
	}{
		{filepath.Join(testDir, "subdir1"), "file1.c"},
		{filepath.Join(testDir, "subdir1"), "file2.cpp"},
		{filepath.Join(testDir, "subdir2"), "file3.h"},
		{filepath.Join(testDir, "subdir2/subsubdir1"), "file4.hpp"},
		{filepath.Join(testDir, "subdir2/subsubdir2"), "file5.cpp"},
	}

	for _, file := range files {
		if err := os.WriteFile(filepath.Join(file.path, file.name), []byte("test"), 0644); err != nil {
			return "", err
		}
	}

	return testDir, nil
}

func TestSearchFiles(t *testing.T) {
	testDir, err := setupTestDir()
	if err != nil {
		t.Fatalf("Failed to set up test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	tests := []struct {
		name       string
		testDir    string
		extensions []string
		expected   []string
		expectErr  error
	}{
		{
			name:       "find C files",
			testDir:    testDir,
			extensions: []string{".c"},
			expected:   []string{filepath.Join(testDir, "subdir1", "file1.c")},
			expectErr:  nil,
		},
		{
			name:       "find header files",
			testDir:    testDir,
			extensions: []string{".h"},
			expected:   []string{filepath.Join(testDir, "subdir2", "file3.h")},
			expectErr:  nil,
		},
		{
			name:       "find multiple extensions",
			testDir:    testDir,
			extensions: []string{".h", ".hpp"},
			expected: []string{
				filepath.Join(testDir, "subdir2", "file3.h"),
				filepath.Join(testDir, "subdir2", "subsubdir1", "file4.hpp"),
			},
			expectErr: nil,
		},
		{
			name:       "no matching files",
			testDir:    testDir,
			extensions: []string{".txt"},
			expected:   []string{},
			expectErr:  nil,
		},
		{
			name:       "empty search extensions",
			testDir:    testDir,
			extensions: []string{},
			expected:   []string{},
			expectErr:  ErrNoExtensions,
		},
		{
			name:       "non-existent directory",
			testDir:    "nonexist",
			extensions: []string{".c"},
			expected:   []string{},
			expectErr:  ErrDirectoryNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := SearchFiles(tt.testDir, tt.extensions)

			if err != nil {
				assert.Equal(t, tt.expectErr, err)
			}

			// Sort expected and actual results for comparison
			sort.Strings(tt.expected)
			sort.Strings(results)

			// Compare lengths
			assert.Equal(t, len(tt.expected), len(results), "Length mismatch")

			// Compare each result
			for i := range results {
				assert.Equal(t, tt.expected[i], results[i], "Mismatch at index %d", i)
			}
		})
	}
}
