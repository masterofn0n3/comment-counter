package commentcounter

import (
	"bufio"
	"os"
)

type state int

const (
	normal state = iota
	slash
	inline
	block
	blockEnding
	stringLiteral
	rawStringLiteral
	characterLiteral
)

// CountComments counts the number of single-line and multi-line comments in a file
func CountComments(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	inlineComments := 0
	blockComments := 0
	state := normal
	scanner := bufio.NewScanner(file)
	currentLine := 0
	lastBlockCommentLine := -1

	for scanner.Scan() {
		line := scanner.Text()
		currentLine++
		i := 0
		lineLen := len(line)

	inner:
		for i < lineLen {
			char := line[i]

			switch state {
			case normal:
				// If we encounter a slash, we might be at the beginning of a comment
				if char == '/' {
					state = slash
				}

			case slash:
				// If we encounter another slash, we are in a single-line comment
				if char == '/' {
					state = inline
					// If we encounter an asterisk, we are in a multi-line comment
				} else if char == '*' {
					state = block
					// If we encounter anything else, we are not in a comment
				} else {
					state = normal
				}

			case inline:
				break inner

			case block:
				// If we encounter an asterisk, we might be at the end of a multi-line comment
				if char == '*' {
					state = blockEnding
				}

			case blockEnding:
				// If we encounter a slash, we are at the end of a multi-line comment
				// We still need to increment the count before going back to normal state
				if char == '/' {
					state = normal
					// Check if the block comment is on the same line before incrementing the count
					if lastBlockCommentLine != currentLine {
						blockComments++
						lastBlockCommentLine = currentLine
					}
				} else {
					state = block
				}
			}

			i++
		}
		if state == inline {
			inlineComments++
		} else if state == block || state == blockEnding {
			blockComments++
			lastBlockCommentLine = currentLine
		}

		if state == inline {
			state = normal
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return inlineComments, blockComments, nil
}
