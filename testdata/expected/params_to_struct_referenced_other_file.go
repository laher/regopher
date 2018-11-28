package main

import "fmt"

func extractParamReferencedOtherFile(param extractParamReferencedOtherFileParam) {
	fmt.Printf("regopher %s\n", param.a)
}

type extractParamReferencedOtherFileParam struct {
	a, b, c string
}
