package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"compass.com/go-homework/pkg/commentcounter"
)

func main() {
	// dirPath, err := utils.ParseArgs()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	cppCommentCounter := commentcounter.NewCppCommentCounter()
	results, err := commentcounter.RecursiveCount(cppCommentCounter, "./testing")
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Fail to start the server: %s\n", err)
	}

	// utils.PrintResults(results)
}
