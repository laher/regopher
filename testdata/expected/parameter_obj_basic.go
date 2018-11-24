package main

import "fmt"

func basicUnusedParamUnreferenced(param basicUnusedParamUnreferencedParam) {
	fmt.Println("regopher")
}

type basicUnusedParamUnreferencedParam struct {
	a, b, c string
}
