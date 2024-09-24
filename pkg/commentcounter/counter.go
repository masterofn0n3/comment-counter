package commentcounter

import (
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
	CountComments(filename string) (int, int, int, error)
	GetExtensions() []string
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

	for i, filePath := range filePaths {
		wg.Add(1)
		go func(i int, filePath string) {
			defer wg.Done()
			total, inlineCount, blockCount, err := counter.CountComments(filePath)
			if err != nil {
				resultsChan <- &CountResult{
					FilePath: filePath,
					Total:    0,
				}
				return
			}
			resultsChan <- &CountResult{
				FilePath:    filePath,
				Total:       total,
				InlineCount: inlineCount,
				BlockCount:  blockCount,
			}
		}(i, filePath)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		for i, filepath := range filePaths {
			if filepath == result.FilePath {
				results[i] = result
				break
			}
		}
	}

	return results, nil
}
