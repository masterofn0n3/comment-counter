package commentcounter

import (
	"errors"
	"reflect"
	"testing"

	"compass.com/go-homework/pkg/filesearch"
)

// MockCounter implements the Counter interface for testing
type MockCounter struct {
	extensions []string
	countFunc  func(string) (*CountResult, error)
}

func (m *MockCounter) CountComments(filename string) (*CountResult, error) {
	return m.countFunc(filename)
}

func (m *MockCounter) GetExtensions() []string {
	return m.extensions
}

// MockFileSearch mocks the filesearch.SearchFiles function
func MockFileSearch(paths []string, err error) func(string, []string) ([]string, error) {
	return func(string, []string) ([]string, error) {
		return paths, err
	}
}

func TestRecursiveCount(t *testing.T) {
	tests := []struct {
		name           string
		counter        Counter
		directoryPath  string
		mockPaths      []string
		mockSearchErr  error
		expectedResult []*CountResult
		expectedErr    string
	}{
		{
			name: "successful return",
			counter: &MockCounter{
				extensions: []string{".go"},
				countFunc: func(filePath string) (*CountResult, error) {
					switch filePath {
					case "file1.go":
						return &CountResult{FilePath: "file1.go", Total: 10, InlineCount: 5, BlockCount: 5}, nil
					case "file2.go":
						return &CountResult{FilePath: "file2.go", Total: 20, InlineCount: 10, BlockCount: 10}, nil
					default:
						return &CountResult{FilePath: "", Total: 0, InlineCount: 0, BlockCount: 0}, errors.New("unexpected file")
					}
				},
			},
			directoryPath: "/test",
			mockPaths:     []string{"file1.go", "file2.go"},
			mockSearchErr: nil,
			expectedResult: []*CountResult{
				{FilePath: "file1.go", Total: 10, InlineCount: 5, BlockCount: 5},
				{FilePath: "file2.go", Total: 20, InlineCount: 10, BlockCount: 10},
			},
			expectedErr: "",
		},
		{
			name: "SearchFiles func return error",
			counter: &MockCounter{
				extensions: []string{".go"},
			},
			directoryPath:  "/test",
			mockPaths:      nil,
			mockSearchErr:  errors.New("search error"),
			expectedResult: nil,
			expectedErr:    "search error",
		},
		{
			name: "CountComments func return error",
			counter: &MockCounter{
				extensions: []string{".go"},
				countFunc: func(filePath string) (*CountResult, error) {
					return &CountResult{FilePath: "", Total: 0, InlineCount: 0, BlockCount: 0}, errors.New("some error")
				},
			},
			directoryPath:  "/test",
			mockPaths:      []string{"file1.go"},
			mockSearchErr:  nil,
			expectedResult: nil,
			expectedErr:    "error processing file file1.go: some error",
		},
		{
			name: "empty directory",
			counter: &MockCounter{
				extensions: []string{".go"},
			},
			directoryPath:  "/test",
			mockPaths:      []string{},
			mockSearchErr:  nil,
			expectedResult: []*CountResult{},
			expectedErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalSearchFiles := filesearch.SearchFiles
			filesearch.SearchFiles = MockFileSearch(tt.mockPaths, tt.mockSearchErr)
			defer func() { filesearch.SearchFiles = originalSearchFiles }()

			result, err := RecursiveCount(tt.counter, tt.directoryPath)

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error %q, but got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("Expected error %q, but got %q", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
