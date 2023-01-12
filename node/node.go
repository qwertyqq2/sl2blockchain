package node

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/crypto"
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
	addr      string
}

const (
	MaxTx = 2
)

func NewNode(user *blockchain.User, addr string) *Node {
	pkstr := crypto.StringPrivate(user.Private())
	err := writeFile("private"+addr, pkstr)
	if err != nil {
		panic("cant write public")
	}
	pybstr := crypto.StringPublic(&user.Private().PublicKey)
	err = writeFile("public"+addr, pybstr)
	if err != nil {
		panic("cant write public")
	}
	addrs := make([]string, 10)
	data, err := readFile(neighborsFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &addrs)
	if err != nil {
		panic(err)
	}
	node := &Node{
		user: user,
	}
	node.neighbors = make(map[string]bool)
	for _, a := range addrs {
		node.neighbors[a] = true //добавить пинг
	}
	var addresses []string
	for addr := range node.neighbors {
		addresses = append(addresses, addr)
	}
	node.hub = hub.NewHub(addresses, dbname, node.addr)
	node.addr = addr
	return node
}
func (n *Node) Public() string {
	return n.user.Public()
}
func LoadNode(addr string) (*Node, error) {
	pk, err := readFile("private" + addr)
	if err != nil {
		return nil, err
	}
	addrs := make([]string, 10)
	data, err := readFile(neighborsFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &addrs)
	if err != nil {
		return nil, err
	}
	user := blockchain.ParseUser(pk)
	n := &Node{
		user: user,
	}
	n.neighbors = make(map[string]bool)
	for _, a := range addrs {
		n.neighbors[a] = true //добавить пинг
	}
	var addresses []string
	for addr := range n.neighbors {
		addresses = append(addresses, addr)
	}
	n.addr = addr
	n.hub = hub.NewHub(addresses, dbname, n.addr)
	return n, nil
}
func (n *Node) NewChain() error {
	gen, err := blockchain.NewBlockchain(dbname+n.addr, n.user.Public())
	if err != nil {
		return err
	}
	genstr, err := blockchain.SerializeBlock(gen)
	if err != nil {
		return err
	}
	pkg := &network.Package{
		Option: OptGenesisBlock,
		Data:   genstr,
	}
	hub.InsertIntoHub(pkg, n.hub)
	return nil
}
func (n *Node) PrintBc() error {
	bc, err := blockchain.Load(dbname + n.addr)
	if err != nil {
		return err
	}
	bc.PrintBlockchain()
	return nil
}
func (n *Node) NewBlockTesting(receiver string, val uint64) error {
	bc, err := blockchain.Load(dbname + n.addr)
	if err != nil {
		return err
	}
	block := blockchain.NewBlock(n.user.Public(), bc.LastBlock().Hash())
	tx1, err := blockchain.NewTransaction(n.user, bc.LastBlock().Hash(), receiver, val)
	if err != nil {
		return err
	}
	err = block.InsertTx(bc, tx1)
	if err != nil {
		return err
	}
	err = block.Accept(bc, n.user, nil)
	if err != nil {
		return err
	}
	err = n.insertBlock(block, bc)
	if err != nil {
		return err
	}
	return nil
}
func (n *Node) insertBlock(block *blockchain.Block, bc *blockchain.Blockchain) error {
	err := bc.InsertBlock(block)
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
func (n *Node) AddBlock(block *blockchain.Block, bc *blockchain.Blockchain) error {
	if err := block.IsValid(bc); err != nil {
		if !bytes.Equal(block.PrevHash, bc.LastBlock().CurHash) {
			if bc.BlockInChain(block) {
				return ErrBlockAlreadyExist
			}
			pkg := &network.Package{
				Option: OptGetBlocks,
				Data:   crypto.Base64Encode(bc.LastBlock().CurHash),
			}
			hub.InsertIntoHub(pkg, n.hub)
			return nil
		}
		return err
	}
	err := n.insertBlock(block, bc)
	if err != nil {
		return err
	}
	return err
}
func (n *Node) CreateTx(receiver string, value uint64) error {
	bc, err := blockchain.Load(dbname + n.addr)
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
	errChan := make(chan error, 5)
	txChan := make(chan []*blockchain.Transaction, MaxTx)
	n.hub.ListenHub(errChan, txChan, MaxTx)
	for {
		select {
		case <-errChan:

		case txs := <-txChan:
			go func() {
				bc, err := blockchain.Load(dbname + n.addr)
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
				err = n.insertBlock(block, bc)
				if err != nil {
					log.Fatal(err)
				}
			}()
		case pkg := <-n.hub.PendingResp:
			go func() {
				err := handlePendingPack(pkg, n)
				if err != nil {
					errChan <- err
				}
			}()
		}
	}
}
func handlePendingPack(pkg *network.Package, n *Node) error {
	switch pkg.Option {
	case OptGetBlocks:
		blocks, err := blockchain.DeserializeBlocks(pkg.Data)
		if err != nil {
			return err
		}
		bc, err := blockchain.Load(dbname + n.addr)
		if err != nil {
			return err
		}
		beginHash := blocks[0].PrevHash
		if !bytes.Equal(beginHash, bc.LastBlock().CurHash) {
			return err
		}
		for _, b := range blocks {
			if err := b.IsValid(bc); err != nil {
				return ErrNotValidBlock
			}
			err = n.insertBlock(b, bc)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
