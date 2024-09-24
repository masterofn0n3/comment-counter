package commentcounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountComments(t *testing.T) {
	tests := []struct {
		name                string
		filePath            string
		expectedTotalLines  int
		expectedInlineCount int
		expectedBlockCount  int
		expectedError       error
	}{
		{
			name:                "file only have inline comment",
			filePath:            "../../test/testdata/only_inline.cpp",
			expectedTotalLines:  10,
			expectedInlineCount: 5,
			expectedBlockCount:  0,
			expectedError:       nil,
		},
		{
			name:                "file only have block comment",
			filePath:            "../../test/testdata/only_block.cpp",
			expectedTotalLines:  10,
			expectedInlineCount: 0,
			expectedBlockCount:  5,
			expectedError:       nil,
		},
		{
			name:                "file have literal double quote",
			filePath:            "../../test/testdata/literal_double_quote.cpp",
			expectedTotalLines:  2,
			expectedInlineCount: 1,
			expectedBlockCount:  0,
			expectedError:       nil,
		},
		{
			name:                "file have literal single quote",
			filePath:            "../../test/testdata/literal_single_quote.cpp",
			expectedTotalLines:  2,
			expectedInlineCount: 1,
			expectedBlockCount:  0,
			expectedError:       nil,
		},
		{
			name:                "file have multi inline with back slash",
			filePath:            "../../test/testdata/multi_line_inline_with_back_slash.cpp",
			expectedTotalLines:  2,
			expectedInlineCount: 2,
			expectedBlockCount:  0,
			expectedError:       nil,
		},
		{
			name:                "file have raw string literal with delimiter",
			filePath:            "../../test/testdata/raw_string_literal_with_delimiter.cpp",
			expectedTotalLines:  9,
			expectedInlineCount: 0,
			expectedBlockCount:  0,
			expectedError:       nil,
		},
		{
			name:                "file have raw string literal with no delimeter",
			filePath:            "../../test/testdata/raw_string_literal.cpp",
			expectedTotalLines:  5,
			expectedInlineCount: 0,
			expectedBlockCount:  0,
			expectedError:       nil,
		},
		{
			name:                "file have 2 block comment on the same line",
			filePath:            "../../test/testdata/two_block_on_same_line.cpp",
			expectedTotalLines:  8,
			expectedInlineCount: 0,
			expectedBlockCount:  4,
			expectedError:       nil,
		},
		{
			name:                "Test special cases",
			filePath:            "../../testing/cpp/special_cases.cpp",
			expectedTotalLines:  62,
			expectedInlineCount: 6,
			expectedBlockCount:  34,
			expectedError:       nil,
		},
		{
			name:                "Non-existent file",
			filePath:            "../../testing/cpp/lib_json/non_existent.cpp",
			expectedTotalLines:  0,
			expectedInlineCount: 0,
			expectedBlockCount:  0,
			expectedError:       assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cppCounter := &CppCommentCounter{}
			totalLines, inlineCount, blockCount, err := cppCounter.CountComments(tt.filePath)
			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedTotalLines, totalLines, "wrong total lines")
			assert.Equal(t, tt.expectedInlineCount, inlineCount, "wrong inline comment count")
			assert.Equal(t, tt.expectedBlockCount, blockCount, "wrong block comment count")
		})
	}
}
