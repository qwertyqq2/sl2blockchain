package main

import (
	"log"

	"github.com/qwertyqq2/sl2blockchain/node"
)

func main() {
	node := node.NewNode()
	err := node.NewChain()
	if err != nil {
		log.Fatal(err)
	}
}
