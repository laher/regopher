package main

import "fmt"

func extractParamUsed(param extractParamUsedParam) {
	fmt.Printf("regopher %s\n", param.a)
}

type extractParamUsedParam struct {
	a, b, c string
}
