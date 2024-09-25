package commentcounter

import (
	"fmt"
	"sync"

	"compass.com/go-homework/pkg/filesearch"
)

type CountResult struct {
	FilePath    string
	Total       int
	InlineCount int
	BlockCount  int
}

type Counter interface {
	CountComments(filename string) (*CountResult, error)
	GetExtensions() []string
}

type FileProcessingError struct {
	FilePath string
	Err      error
}

func (e *FileProcessingError) Error() string {
	return fmt.Sprintf("error processing file %s: %v", e.FilePath, e.Err)
}

// RecursiveCount counts comments in all files in a directory and its subdirectories
func RecursiveCount(counter Counter, directoryPath string) ([]*CountResult, error) {
	filePaths, err := filesearch.SearchFiles(directoryPath, counter.GetExtensions())
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make([]*CountResult, len(filePaths))
	resultsChan := make(chan *CountResult, len(filePaths))
	errChan := make(chan error, 1)

	for i, filePath := range filePaths {
		wg.Add(1)
		go func(i int, filePath string) {
			defer wg.Done()
			countResult, err := counter.CountComments(filePath)
			if err != nil {
				errChan <- &FileProcessingError{FilePath: filePath, Err: err}
				return
			}
			resultsChan <- countResult
		}(i, filePath)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for i := 0; i < len(filePaths); i++ {
		select {
		case err := <-errChan:
			return nil, err
		case result := <-resultsChan:
			for i, filepath := range filePaths {
				if filepath == result.FilePath {
					results[i] = result
					break
				}
			}
		}
	}

	return results, nil
}
