package main

import (
	"fmt"
	"log"
)

// A Product contains metadata about a product for sale.

func main() {
	store, err := NewStore()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("made it here bitch")
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":4000", store)
	server.Run()
}
