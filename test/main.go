package main

import "fmt"

type Signal int

func main() {
	SIGNINT := Signal(0x2)
	fmt.Println(2 == 0x2)
	fmt.Printf("%d ->  %T", SIGNINT, SIGNINT)
}
