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
	characterLiteral

	// States for raw string literals
	preRawStringLiteral
	confirmedRawStringLiteral
	rawStringLiteral
	rawStringLiteralEnding
	preRawStringLiteralDelimiter
	rawStringLiteralDelimiter
	rawStringLiteralDelimiterEnding
)

// CountComments counts the number of single-line and multi-line comments in a file
func CountComments(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	inlineComments := 0
	inlineContinuation := false
	blockComments := 0
	state := normal
	reader := bufio.NewReader(file)
	currentLine := 0
	lastBlockCommentLine := -1
	var delimiter string

	for {
		char, _, err := reader.ReadRune()
		charS := string(char)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return 0, 0, err
		}
		switch state {
		case normal:
			// If we encounter a slash, we might be at the beginning of a comment
			if char == '/' {
				state = slash
			} else if char == '"' {
				state = stringLiteral
			} else if char == 'R' {
				state = preRawStringLiteral
			}
		case stringLiteral:
			if char == '"' {
				state = normal
			}
		case preRawStringLiteral:
			if char == '"' {
				state = confirmedRawStringLiteral
			} else {
				state = normal
			}
		case confirmedRawStringLiteral:
			if char == '(' {
				state = rawStringLiteral
			} else {
				state = preRawStringLiteralDelimiter
				delimiter += charS
			}
		case rawStringLiteral:
			if char == ')' {
				state = rawStringLiteralEnding
			}
		case rawStringLiteralEnding:
			if char == '"' {
				state = normal
			} else {
				state = rawStringLiteral
			}
		case preRawStringLiteralDelimiter:
			if char == '(' {
				state = rawStringLiteralDelimiter
			} else {
				delimiter += charS
			}
		case rawStringLiteralDelimiter:
			if char == ')' {
				peeked, err := reader.Peek(len(delimiter))
				if err != nil {
					return 0, 0, err
				}
				if (string(peeked)) == delimiter {
					state = rawStringLiteralDelimiterEnding
					_, err := reader.Discard(len(delimiter))
					if err != nil {
						return 0, 0, err
					}
				}
			}
		case rawStringLiteralDelimiterEnding:
			if char == '"' {
				state = normal
			} else {
				state = rawStringLiteralDelimiter
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
			if char == '\\' {
				inlineContinuation = true
				break
			}
			if inlineContinuation {
				if char == '\n' {
					inlineComments++
				}
				inlineContinuation = false
			}

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
		if char == '\n' {
			if state == inline {
				state = normal
				inlineComments++
			} else if state == block || state == blockEnding {
				blockComments++
				lastBlockCommentLine = currentLine
			}
			currentLine++
		}
	}

	return inlineComments, blockComments, nil
}
