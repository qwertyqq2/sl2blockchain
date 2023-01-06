package node

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
	"github.com/qwertyqq2/sl2blockchain/node/hub"
)

var (
	dbname        = "blockchain.db"
	neighborsFile = "addr.json"
)

type Node struct {
	user      *blockchain.User
	neighbors map[string]bool
	hub       *hub.Hub
}

func NewNode(user *blockchain.User) *Node {
	err := writeFile("node-params", user.Public())
	if err != nil {
		panic("cant write public")
	}
	addr := make([]string, 10)
	data, err := readFile(neighborsFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &addr)
	if err != nil {
		panic(err)
	}
	node := &Node{
		user: user,
	}
	node.neighbors = make(map[string]bool)
	for _, a := range addr {
		node.neighbors[a] = true //добавить пинг
	}
	var addresses []string
	for addr := range node.neighbors {
		addresses = append(addresses, addr)
	}
	node.hub = hub.NewHub(addresses, dbname)
	return node
}

func (n *Node) NewChain() error {
	gen, err := blockchain.NewBlockchain(dbname, n.user.Public())
	if err != nil {
		return err
	}
	genstr, err := blockchain.SerializeBlock(gen)
	if err != nil {
		return err
	}
	pkg := &network.Package{
		Option: OptAddBlock,
		Data:   genstr,
	}
	hub.InsertIntoHub(pkg, n.hub)
	return nil
}

func (n *Node) insertBlock(block *blockchain.Block) error {
	bc, err := blockchain.Load(dbname)
	if err != nil {
		return err
	}
	err = bc.InsertBlock(block)
	if err != nil {
		return err
	}
	blockstr, err := blockchain.SerializeBlock(block)
	if err != nil {
		return err
	}
	pkg := &network.Package{
		Option: OptAddBlock,
		Data:   blockstr,
	}
	hub.InsertIntoHub(pkg, n.hub)
	return nil
}

func (n *Node) CreateTx(receiver string, value uint64) error {
	bc, err := blockchain.Load(dbname)
	if err != nil {
		return err
	}
	lasthash := bc.LastBlock().CurHash
	tx, err := blockchain.NewTransaction(n.user, lasthash, receiver, value)
	if err != nil {
		return err
	}
	txser, err := blockchain.SerializeTX(tx)
	if err != nil {
		return err
	}
	pkg := &network.Package{
		Option: OptNewTx,
		Data:   txser,
	}
	hub.InsertIntoHub(pkg, n.hub)
	return nil
}

func (n *Node) HubCheck() {
	errChan := make(chan error)
	txChan := make(chan []*blockchain.Transaction, 3)
	n.hub.ListenHub(errChan, txChan)
	for {
		select {
		case err := <-errChan:
			log.Fatal(err)
		case txs := <-txChan:
			bc, err := blockchain.Load(dbname)
			if err != nil {
				log.Fatal(err)
			}
			block := blockchain.NewBlock(n.user.Public(), bc.LastBlock().CurHash)
			for _, tx := range txs {
				err := block.InsertTx(bc, tx)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = block.Accept(bc, n.user, nil)
			if err != nil {
				log.Fatal(err)
			}
			err = n.insertBlock(block)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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
