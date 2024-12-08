package main

import "fmt"

func main() {
	done := make(chan struct{}) // channel act as a signal for synchronizing goroutines
	go pipelines(done)
	<-done // waits for a task to complete
}

func pipelines(done chan struct{}) {
	numbers := make(chan int)
	squares := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			numbers <- i
		}
		close(numbers)
	}()
	go func() {
		for num := range numbers {
			squares <- num * num
		}
		close(squares)
	}()

	for square := range squares {
		fmt.Println(square)
	}
	done <- struct{}{} // signal completion
}

func worker() {
	
}