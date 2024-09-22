package commentcounter

import (
	"bufio"
	"os"
)

type State int

const (
	NORMAL State = iota
	SLASH
	SINGLE_LINE_COMMENT
	MULTI_LINE_COMMENT
	MULTI_LINE_TERMINATING
)

// CountComments counts the number of single-line and multi-line comments in a file
func CountComments(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	singleLineComments := 0
	multiLineComments := 0
	state := NORMAL
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		i := 0
		lineLen := len(line)

	inner:
		for i < lineLen {
			char := line[i]

			switch state {
			case NORMAL:
				if char == '/' {
					state = SLASH
				}

			case SLASH:
				if char == '/' {
					state = SINGLE_LINE_COMMENT
					singleLineComments++
				} else if char == '*' {
					state = MULTI_LINE_COMMENT
					multiLineComments++
				} else {
					state = NORMAL
				}

			case SINGLE_LINE_COMMENT:
				break inner

			case MULTI_LINE_COMMENT:
				if char == '*' {
					state = MULTI_LINE_TERMINATING
				}
				if i == 0 {
					multiLineComments++
				}

			case MULTI_LINE_TERMINATING:
				if char == '/' {
					state = NORMAL
				} else if char != '*' {
					state = MULTI_LINE_COMMENT
				}
			}

			i++
		}

		if state == SINGLE_LINE_COMMENT {
			state = NORMAL
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return singleLineComments, multiLineComments, nil
}
