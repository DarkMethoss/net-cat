package main

import (
	"fmt"
	"time"
)

func main() {
	skip := make(chan string)

	// Goroutine to send a value to the channel
	go func() {
		time.Sleep(time.Second*3)
		skip <- "done"
	}()

	fmt.Println("1")
	time.Sleep(time.Second*1)
	fmt.Println("2")
	time.Sleep(time.Second*1)
	fmt.Println("3")
	fmt.Println("4")
	time.Sleep(time.Second*1)
	fmt.Println("5")
	time.Sleep(time.Second*1)
	fmt.Println("6")
	time.Sleep(time.Second*1)
	fmt.Println("7")
	time.Sleep(time.Second*1)
	fmt.Println("8")
	time.Sleep(time.Second*1)
	fmt.Println("9")
	time.Sleep(time.Second*1)
	fmt.Println("10")

	// Receive from the channel
	<-skip
	fmt.Println("Exiting program without deadlock.")
}

//! What This Code Does Not Do
// It doesn't "skip" or control the output of any specific lines; all fmt.Println statements are 
// executed in order. The channel usage here doesn't affect the output flow.