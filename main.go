package main

import (
	"fmt"
	"rsc.io/quote"
	"github.com/memory1/gomodtest"

)

//CompareLength ...
func CompareLength(l1, l2 float32) bool {
	return l1 == l2
}

func main() {
	fmt.Println(quote.Hello()) 
	eq := CompareLength(1,2)
	if eq {
		fmt.Println("equal")
	} else {
		fmt.Println("false")
	}

	eq2 := gomodtest.Mile(3).Equal(gomodtest.Mile(3))
	if eq2 {
		fmt.Println("equal")
	} else {
		fmt.Println("false")
	}
}