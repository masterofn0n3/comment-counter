package main

import (
	"fmt"
)

func main() {
	// Sample input: Replace this with the content of your file
	code := `
    main()
    {
        printf(
            "// This is not a comment"
            "// This is not a comment\
            // This is not a comment" // Comment: inline
            "// This is not a comment \" not a comment as well"
            "/* This is not a comment */"
            "/* This is not a comment \
            This is not a comment */"
            "This is not a comment */"
            "// This line has a trailing comment" // Comment: \
            inline (the line-break makes this one counted as 2 lines)
        );                                        
        /* Comment: block
        Comment: block
        Comment: block
        */ int x = 0;                                
        // Comment: block, inline
        /*
        Comment: block
        */ int y = 1;                                
        /*
        Comment: block */

        /* NOTE:
        The next block follows the raw string literals syntax, which is a special case
        for string parsing.
        In short, the format is R"[delimiter](...)[delimiter]", where the content
        between the quotes (") are not comments.
        The delimiters before and after the parentheses must match.

        You can find more details here:
        https://en.cppreference.com/w/cpp/language/string_literal
        */

        const char * vogon_poem = R"V0G0N(
        // NOT A COMMENT:             O freddled gruntbuggly thy micturations are to me
        // NOT A COMMENT:                As plured gabbleblochits on a lurgid bee.
        // NOT A COMMENT:             Groop, I implore thee my foonting turlingdromes.
        /* NOT A COMMENT:          And hooptiously drangle me with crinkly bindlewurdles,
         * NOT A COMMENT: Or I will rend thee in the gobberwarts with my blurlecruncheon, see if I don't.
         */<<NOT A COMMENT
                        (by Prostetnic Vogon Jeltz; see p. 56/57)
        )V0G0N";
    }

    /* Comment: block
    Comment: block \
    /* Comment: block // This is NOT an inline comment
    /* Comment: block
    /* Comment: block
    Comment: block */ // Comment: inline // Not another inline comment
    //* Comment: inline
    `

	// Define states
	const (
		Normal = iota
		Slash
		InlineComment
		BlockComment
		BlockCommentEnd
		StringLiteral
		RawStringLiteral
		CharacterLiteral
	)

	state := Normal
	inlineCommentCount := 0
	blockCommentCount := 0
	delimiter := rune(0)

	for i := 0; i < len(code); i++ {
		ch := rune(code[i])

		switch state {
		case Normal:
			if ch == '/' {
				state = Slash
			} else if ch == '"' {
				state = StringLiteral
			} else if ch == 'R' && i+1 < len(code) && code[i+1] == '"' {
				state = RawStringLiteral
				i++ // Skip the opening quote
			} else if ch == '\'' {
				state = CharacterLiteral
			}

		case Slash:
			if ch == '/' {
				state = InlineComment
				inlineCommentCount++
			} else if ch == '*' {
				state = BlockComment
				blockCommentCount++
			} else {
				state = Normal
			}

		case InlineComment:
			if ch == '\n' {
				state = Normal
			}

		case BlockComment:
			if ch == '*' {
				state = BlockCommentEnd
			}

		case BlockCommentEnd:
			if ch == '/' {
				state = Normal
			} else {
				state = BlockComment
			}

		case StringLiteral:
			if ch == '\\' && i+1 < len(code) {
				i++ // Skip escaped character
			} else if ch == '"' {
				state = Normal
			}

		case RawStringLiteral:
			if delimiter == 0 {
				delimiter = ch
			} else if ch == delimiter && i+1 < len(code) && code[i+1] == '"' {
				state = Normal
				i++ // Skip the closing quote
				delimiter = 0
			}

		case CharacterLiteral:
			if ch == '\\' && i+1 < len(code) {
				i++ // Skip escaped character
			} else if ch == '\'' {
				state = Normal
			}
		}
	}

	fmt.Printf("Inline Comments: %d\n", inlineCommentCount)
	fmt.Printf("Block Comments: %d\n", blockCommentCount)
}
