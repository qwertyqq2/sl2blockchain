package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
		addr = os.Args[2]
		n = node.NewNode(user, addr)
	case "loadnode":
		addr = os.Args[2]
		n_, err := node.LoadNode(addr)
		if err != nil {
			panic(err)
		}
		n = n_
		fmt.Println("pubkey: ", n.Public())
	}
}

func main() {
	go n.HubCheck()
	handle := node.NewHandleNode(n)
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
			}
		case "newb":
			if len(splited) < 2 {
				fmt.Println("not enouhg args")
				continue
			}
			receiver := splited[1]
			val, err := strconv.ParseUint(splited[2], 10, 64)
			err = n.NewBlockTesting(receiver, val)
			if err != nil {
				fmt.Println(err)
			}
		case "ctx":
			if len(splited) < 2 {
				fmt.Println("not enouhg args")
				continue
			}
			receiver := splited[1]
			val, err := strconv.ParseUint(splited[2], 10, 64)
			if err != nil {
				log.Println(err)
			}
			err = n.CreateTx(receiver, val)
			if err != nil {
				fmt.Println(err)
			}
		case "printbc":
			err := n.PrintBc()
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("undefined")
		}
	}
}

func input() string {
	fmt.Printf("->")
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic("error reader")
	}
	return msg
}
