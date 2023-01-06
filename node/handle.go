package node

import (
	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
)

var (
	OptGenesisBlock = 200
	OptAddBlock     = 201
	OptGetBlock     = 202
	OptLastHash     = 203
	OptGetBalance   = 204
	OptGetSizeChain = 205
	OptNewTx        = 206
)

type Handle struct {
	node *Node
}

func NewHandleNode() *Handle {
	return &Handle{}
}

func (h *Handle) handleServer(conn network.Conn, pack *network.Package) {
	network.Handle(OptGenesisBlock, conn, pack, newChain)
	network.Handle(OptAddBlock, conn, pack, addBlock)
	// network.Handle(OptGetBlock, conn, pack, getBlock)
	// network.Handle(OptLastHash, conn, pack, getLastHash)
	// network.Handle(OptGetBalance, conn, pack, getBalance)
	// network.Handle(OptGetSizeChain, conn, pack, getChainSize)
}

func (h *Handle) Listen(addr string) {
	network.Listen(addr, h.handleServer)
}

func newChain(pkg *network.Package) string {
	genesisstr := pkg.Data
	block, err := blockchain.DeserializeBlock(genesisstr)
	if err != nil {
		return "incorrect block"
	}
	_, err = blockchain.NewBlockchainWithGenesis(block, dbname, block.Miner)
	if err != nil {
		return "some..."
	}
	return "yes"
}

func addBlock(pkg *network.Package) string {
	blockstr := pkg.Data
	block, err := blockchain.DeserializeBlock(blockstr)
	if err != nil {
		return "incorrectBlock"
	}
	bc, err := blockchain.Load(dbname)
	if err != nil {
		return "notBlockchain"
	}
	if !block.IsValid(bc) {
		return "notValid"
	}
	err = bc.InsertBlock(block)
	if err != nil {
		return "cantInsertBlock"
	}
	return "yes"
}
