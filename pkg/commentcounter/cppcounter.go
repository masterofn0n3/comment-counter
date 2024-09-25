package commentcounter

import (
	"bufio"
	"os"
)

const (
	CFileExtension     = ".c"
	CPPFileExtension   = ".cpp"
	CHeaderExtension   = ".h"
	CPPHeaderExtension = ".hpp"
)

type CppCommentCounter struct {
	extensions []string
}

func NewCppCommentCounter() *CppCommentCounter {
	return &CppCommentCounter{
		extensions: []string{CFileExtension, CPPFileExtension, CHeaderExtension, CPPHeaderExtension},
	}
}

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

// CountComments counts the number of total line, single-line and multi-line comments in a file
func (cppComCounter *CppCommentCounter) CountComments(filename string) (*CountResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &CountResult{FilePath: filename, Total: 0, InlineCount: 0, BlockCount: 0}, err
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
			return &CountResult{FilePath: filename, Total: 0, InlineCount: 0, BlockCount: 0}, err
		}
		switch state {
		case normal:
			if char == '/' {
				state = slash
			} else if char == '"' {
				state = stringLiteral
			} else if char == 'R' {
				state = preRawStringLiteral
			} else if char == '\'' {
				state = characterLiteral
			}
		case characterLiteral:
			if char == '\\' {
				nextChar, _, err := reader.ReadRune()
				if err != nil {
					return &CountResult{FilePath: filename, Total: 0, InlineCount: 0, BlockCount: 0}, err
				}
				if nextChar == '\'' || nextChar == '\\' {
					continue
				} else {
					reader.UnreadRune()
				}
			} else if char == '\'' {
				state = normal
			}
		case stringLiteral:
			if char == '\\' {
				nextChar, _, err := reader.ReadRune()
				if err != nil {
					return &CountResult{FilePath: filename, Total: 0, InlineCount: 0, BlockCount: 0}, err
				}
				if nextChar == '"' || nextChar == '\\' {
					continue
				} else {
					reader.UnreadRune()
				}
			} else if char == '"' {
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
					return &CountResult{FilePath: filename, Total: 0, InlineCount: 0, BlockCount: 0}, err
				}
				// Peeked ahead to see if the delimiter is coming up
				if (string(peeked)) == delimiter {
					state = rawStringLiteralDelimiterEnding
					_, err := reader.Discard(len(delimiter))
					if err != nil {
						return &CountResult{FilePath: filename, Total: 0, InlineCount: 0, BlockCount: 0}, err
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
			if char == '/' {
				state = inline
			} else if char == '*' {
				state = block
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
			if char == '*' {
				state = blockEnding
			}

		case blockEnding:
			if char == '/' {
				state = normal
				// Only increment block comments if the last block comment was on a different line
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
				if lastBlockCommentLine != currentLine {
					blockComments++
					lastBlockCommentLine = currentLine
				}
			}
			currentLine++
		}
	}

	return &CountResult{FilePath: filename, Total: currentLine, InlineCount: inlineComments, BlockCount: blockComments}, nil
}

func (cppComCounter *CppCommentCounter) GetExtensions() []string {
	return cppComCounter.extensions
}
