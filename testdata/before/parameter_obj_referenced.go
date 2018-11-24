package main

import "fmt"

func extractParamReferenced(a, b, c string) {
	fmt.Printf("regopher %s\n", a)
}

func referToParam() {
	extractParamReferenced("1", "2", "3")
}
