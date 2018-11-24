package main

import "fmt"

// foo
func extractParamUnused(param extractParamUnusedParam) {
	fmt.Println("regopher")
}

type extractParamUnusedParam struct {
	a, b, c string
}
