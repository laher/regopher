package main

import "fmt"

type pom struct {
}

// foo
func (p pom) extractParam(param extractParamParam) {
	fmt.Println("regopher")
}

type extractParamParam struct {
	a, b, c string
}
