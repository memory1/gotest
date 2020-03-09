package main

import ("fmt")

func main() {
	naturals := make(chan int)
	squares := make(chan int)
	// counter
	go func() {
		for x := 0; x<10; x++ {
			naturals <- x
			fmt.Printf("x = %d \n",x)
		}
		close(naturals)
	}()
	// squarer
	go func() {
		for {
			x, ok := <-naturals
			if (!ok) {break}
			squares <- x * x
		}
		close(squares)
	}()
	// printer in main goroutine
	for {
		x, ok := <-squares
		if (!ok) {break}
		fmt.Printf("squares = %d \n",x)
	}

}