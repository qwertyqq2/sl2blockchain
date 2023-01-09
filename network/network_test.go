package network

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

const (
	ToUpper = 1
	addr    = ":8090"
)

func handleServer(conn Conn, pack *Package) {
	err := Handle(ToUpper, conn, pack, handleToUpper)
	if err != nil {
		fmt.Println(err)
	}
}

func handleToUpper(pack *Package) string {
	return strings.ToUpper(pack.Data)
}

func TestSend(t *testing.T) {
	_, err := Listen(addr, handleServer)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1000 * time.Millisecond)

	res := Send(addr, &Package{
		Option: ToUpper,
		Data:   "Hello world",
	})
	if err != nil {
		log.Println(err)
	}
	t.Log(res.Data)
}
