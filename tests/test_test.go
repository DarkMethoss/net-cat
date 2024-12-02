package tests

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"testing"
// )

// func TestListener(t *testing.T) {
// 	listner, err := net.Listen("tcp", "localhost:8080")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer func() {
// 		err = listner.Close()
// 		if err != nil {
// 			log.Println("error closing the listner server")
// 		}
// 	}()
// 	fmt.Printf("bound to %q\n", listner.Addr())
// }
