package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/node"
)

var (
	n    *node.Node
	addr string
)

func init() {
	if len(os.Args) < 2 {
		panic("<2 args")
	}
	arg := os.Args[1]
	switch arg {
	case "newnode":
		user := blockchain.NewUser()
		n = node.NewNode(user)
		addr = os.Args[2]
	case "loadnode":
		user := blockchain.ParseUser("qwe") /////parsing
		n = node.NewNode(user)
		addr = os.Args[2]
	}
}

func main() {
	go n.HubCheck()
	handle := node.NewHandleNode()
	go handle.Listen(addr)
	for {
		msgn := input()
		msg := strings.Replace(msgn, "\n", " ", -1)
		splited := strings.Split(msg, " ")
		switch splited[0] {
		case "newchain":
			err := n.NewChain()
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func input() string {
	fmt.Print("->")
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic("error reader")
	}
	return msg
}
