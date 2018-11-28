package main

import "fmt"

func extractParamCollisionSameType(a, b, c string) {
	fmt.Printf("regopher %s\n", a)
	x := []string{"d", "e"}
	for _, a := range x {
		fmt.Printf("regopher %s\n", a)
	}
}
