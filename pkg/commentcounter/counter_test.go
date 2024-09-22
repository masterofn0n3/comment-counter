package commentcounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountComments(t *testing.T) {
	tests := []struct {
		name                string
		filePath            string
		expectedInlineCount int
		expectedBlockCount  int
		expectedError       error
	}{
		// {
		// 	name:                "Test json_reader.cpp having only inline comments",
		// 	filePath:            "../../testing/cpp/lib_json/json_reader.cpp",
		// 	expectedInlineCount: 134,
		// 	expectedBlockCount:  0,
		// 	expectedError:       nil,
		// },
		{
			name:                "Test json_tool.h having both inline and block comments",
			filePath:            "../../testing/cpp/lib_json/json_tool.h",
			expectedInlineCount: 13,
			expectedBlockCount:  19,
			expectedError:       nil,
		},
		// {
		// 	name:                "Test special cases",
		// 	filePath:            "../../testing/cpp/special_cases.cpp",
		// 	expectedInlineCount: 6,
		// 	expectedBlockCount:  34,
		// 	expectedError:       nil,
		// },
		{
			name:                "Non-existent file",
			filePath:            "../../testing/cpp/lib_json/non_existent.cpp",
			expectedInlineCount: 0,
			expectedBlockCount:  0,
			expectedError:       assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inlineCount, blockCount, err := CountComments(tt.filePath)
			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedInlineCount, inlineCount, "wrong inline comment count")
			assert.Equal(t, tt.expectedBlockCount, blockCount, "wrong block comment count")
		})
	}
}
