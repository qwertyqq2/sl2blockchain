package network

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

const (
	TO_UPPER = 1
	ADDRESS1 = ":8090"
)

func handleServer(conn Conn, pack *Package) {
	err := Handle(TO_UPPER, conn, pack, handleToUpper)
	if err != nil {
		fmt.Println(err)
	}
}

func handleToUpper(pack *Package) string {
	return strings.ToUpper(pack.Data)
}

func TestSend(t *testing.T) {
	_, err := Listen(ADDRESS1, handleServer)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1000 * time.Millisecond)

	res := Send(ADDRESS1, &Package{
		Option: TO_UPPER,
		Data:   "Hello world",
	})
	if err != nil {
		log.Println(err)
	}
	t.Log(res.Data)
}
