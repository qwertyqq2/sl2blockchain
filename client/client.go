package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
)

type Client struct {
	user      *blockchain.User
	neighbors []string
}

// /Options
const (
	OptAskBal = 101
	OptGetBal = 102
)

var neighborsFile = "addr.json"

func NewClient(filename string) *Client {
	user := blockchain.NewUser()
	err := writeFile(filename, user.Public())
	if err != nil {
		panic("cant write public")
	}
	addr := make([]string, 2)
	data, err := readFile(neighborsFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &addr)
	if err != nil {
		panic(err)
	}
	return &Client{
		user:      user,
		neighbors: addr,
	}
}

func LoadClient(filename string) *Client {
	pk, err := readFile(filename)
	if err != nil {
		panic("cant load public")
	}
	addr := make([]string, 2)
	data, err := readFile(neighborsFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &addr)
	if err != nil {
		panic(err)
	}
	return &Client{
		user:      blockchain.ParseUser(pk),
		neighbors: addr,
	}
}

func (c *Client) Handle() error {
	for {
		msgn := input()
		msg := strings.Replace(msgn, "\n", " ", -1)
		splited := strings.Split(msg, " ")
		switch splited[0] {
		case "/exit":
			os.Exit(1)
		case "/user":
			if len(splited) < 2 {
				fmt.Println("not enouhg args")
				continue
			}
			switch splited[1] {
			case "private":
				fmt.Println(c.user.Private().N.String())
			case "public":
				fmt.Println(c.user.Public())
			case "addresses":
				fmt.Println(c.neighbors)
			case "balance":
				bal, err := c.getBalance()
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(bal)

			default:
				fmt.Println("undefined cmd")

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

func (c *Client) getBalance() ([]string, error) {
	bal := make([]string, 0)
	for _, addr := range c.neighbors {
		res := network.Send(addr, &network.Package{
			Option: OptGetBal,
			Data:   c.user.Public(),
		})
		if res != nil && res.Option == OptGetBal {
			bal = append(bal, res.Data)
		}
	}
	if bal == nil {
		return nil, ErrNilBalances
	}
	return bal, nil
}

func writeFile(filename string, data string) error {
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
