package main

import "fmt"

func extractParamCollisionSameType(param extractParamCollisionSameTypeParam) {
	fmt.Printf("regopher %s\n", param.a)
	x := []string{"d", "e"}
	for _, a := range x {
		fmt.Printf("regopher %s\n", a)
	}
}

type extractParamCollisionSameTypeParam struct {
	a, b, c string
}
