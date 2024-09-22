package main

import (
	"fmt"

	"compass.com/go-homework/pkg/commentcounter"
)

func main() {
	inline, block, _ := commentcounter.CountComments("testing/cpp/test_lib_json/main.cpp")
	fmt.Println("Inline comments:", inline)
	fmt.Println("Block comments:", block)

}
