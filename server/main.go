package main

import "fmt"

func main() {
	fmt.Println("Hello world")

	server := NewAPIServer(":4000")
	server.Run()
}
