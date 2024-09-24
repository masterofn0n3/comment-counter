package commentcounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCounter struct {
	Files map[string]CountResult
	Err   error
}

func (m *MockCounter) CountComments(filename string) (int, int, int, error) {
	if m.Err != nil {
		return 0, 0, 0, m.Err
	}
	result, exists := m.Files[filename]
	if !exists {
		return 0, 0, 0, nil
	}
	return result.Total, result.InlineCount, result.BlockCount, nil
}

func (m *MockCounter) GetExtensions() []string {
	return []string{".cpp", ".h"}
}

func TestRecursiveCount(t *testing.T) {
	tests := []struct {
		name        string
		mockCounter MockCounter
		directory   string
		expected    []*CountResult
		expectError bool
	}{
		{
			name: "Successful comment counting",
			mockCounter: MockCounter{
				Files: map[string]CountResult{
					"file1.cpp": {Total: 10, InlineCount: 5, BlockCount: 5},
					"file2.cpp": {Total: 15, InlineCount: 10, BlockCount: 5},
				},
			},
			directory: "mock_directory",
			expected: []*CountResult{
				{FilePath: "file1.cpp", Total: 10, InlineCount: 5, BlockCount: 5},
				{FilePath: "file2.cpp", Total: 15, InlineCount: 10, BlockCount: 5},
			},
			expectError: false,
		},
		{
			name: "File not found",
			mockCounter: MockCounter{
				Files: map[string]CountResult{},
			},
			directory: "mock_directory",
			expected: []*CountResult{
				{FilePath: "file1.cpp", Total: 0},
				{FilePath: "file2.cpp", Total: 0},
			},
			expectError: false,
		},
		{
			name: "Error while counting comments",
			mockCounter: MockCounter{
				Err: assert.AnError, // Simulate an error
			},
			directory: "mock_directory",
			expected: []*CountResult{
				{FilePath: "file1.cpp", Total: 0},
				{FilePath: "file2.cpp", Total: 0},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// filePaths := []string{"file1.cpp", "file2.cpp"}

			results, err := RecursiveCount(&tt.mockCounter, tt.directory)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, results)
		})
	}
}
