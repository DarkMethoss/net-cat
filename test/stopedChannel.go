package main

import (
	"fmt"
	"time"
)

var stoped = make(chan struct{})

func main() {
	go test()
	time.Sleep(time.Second * 3)
	close(stoped)
}

func test() {
	for {
		if _stoped
			fmt.Print("stoped")
			return
		default:
			fmt.Println("|")
		}
	}
}
