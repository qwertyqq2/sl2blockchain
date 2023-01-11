package blockchain

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func ranggo() {
	fmt.Println("go rnad")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		slice := make([]byte, 100)
		rand.Read(slice)
		fmt.Println(slice)
	}
}

func TestRand(t *testing.T) {
	ranggo()
	ranggo()
}
