package utils

import (
	"fmt"
	"os"

	"compass.com/go-homework/pkg/commentcounter"
)

// ParseArgs parses command-line arguments and returns the directory path.
func ParseArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("usage: go run . <directory_path>")
	}
	return os.Args[1], nil
}

// PrintResults formats and prints the results.
func PrintResults(results []*commentcounter.CountResult) {
	maxFilePathLength := 0
	for _, result := range results {
		if len(result.FilePath) > maxFilePathLength {
			maxFilePathLength = len(result.FilePath)
		}
	}

	// Create the format string with dynamic width
	formatString := fmt.Sprintf("%%-%ds  total:%%5d  inline:%%5d  block:%%5d\n", maxFilePathLength)

	for _, result := range results {
		fmt.Printf(formatString, result.FilePath, result.Total, result.InlineCount, result.BlockCount)
	}
}
