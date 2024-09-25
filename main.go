package main

import (
	"fmt"

	"compass.com/go-homework/internal/utils"
	"compass.com/go-homework/pkg/commentcounter"
)

func main() {
	dirPath, err := utils.ParseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	cppCommentCounter := commentcounter.NewCppCommentCounter()
	results, err := commentcounter.RecursiveCount(cppCommentCounter, dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.PrintResults(results)
}
