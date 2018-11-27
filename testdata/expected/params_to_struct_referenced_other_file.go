package main

import "fmt"

func extractParamReferenced(param extractParamReferencedParam) {
	fmt.Printf("regopher %s\n", param.a)
}

type extractParamReferencedParam struct {
	a, b, c string
}
