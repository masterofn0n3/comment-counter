package main

import "compass.com/go-homework/pkg/commentcounter"

func main() {
	cppCommentCounter := commentcounter.NewCppCommentCounter()
	result, err := commentcounter.RecursiveCount(cppCommentCounter, "testing/cpp")
	if err != nil {
		println(err)
		return
	}
	for _, res := range result {
		println(res.FilePath, res.Total, res.InlineCount, res.BlockCount)
	}
}
