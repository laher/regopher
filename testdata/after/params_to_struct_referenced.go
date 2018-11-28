package main

import "fmt"

func extractParamReferenced(param extractParamReferencedParam) {
	fmt.Printf("regopher %s\n", param.a)
}

func referToParam() {
	extractParamReferenced(extractParamReferencedParam{"1", "2", "3"})
}

type extractParamReferencedParam struct {
	a, b, c string
}
