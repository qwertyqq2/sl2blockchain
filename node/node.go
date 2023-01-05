package node

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
	"github.com/qwertyqq2/sl2blockchain/node/txpool"
)

var (
	dbname        = "blockchain.db"
	neighborsFile = "addr.json"
)

const (
	OptAskBal   = 101
	OptGetBal   = 102
	OptNewBlock = 200
)

type Node struct {
	user      *blockchain.User
	neighbors map[string]bool
	hub       chan *network.Package
	txpool    *txpool.Pool
}

func NewNode() *Node {
	user := blockchain.NewUser()
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
	node.hub = make(chan *network.Package, len(node.neighbors))
	node.txpool = txpool.NewPool()
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
		Option: OptNewBlock,
		Data:   genstr,
	}
	n.hub <- pkg
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
		Option: OptNewBlock,
		Data:   blockstr,
	}
	n.hub <- pkg
	return nil
}

func (n *Node) ListenHub(errChan chan error) {
	for {
		select {
		case pack := <-n.hub:
			for addr, ok := range n.neighbors {
				if ok {
					resp := network.Send(addr, pack)
					if resp == nil {
						errChan <- ErrNilPackageResp
					}
				}
			}
		}
	}
}

func (n *Node) PoolCheck() {
	for {
		txs, f := n.txpool.GetTxs()
		if f {
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
