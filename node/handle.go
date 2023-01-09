package node

import (
	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
	"github.com/qwertyqq2/sl2blockchain/node/hub"
)

var (
	OptGenesisBlock = 200
	OptAddBlock     = 201
	OptGetBlock     = 202
	OptLastHash     = 203
	OptGetBalance   = 204
	OptGetSizeChain = 205
	OptNewTx        = 206
	Success         = "success"
)

type Handle struct {
	node *Node
}

func NewHandleNode(n *Node) *Handle {
	return &Handle{
		node: n,
	}
}

func (h *Handle) handleServer(conn network.Conn, pack *network.Package) {
	network.Handle(OptGenesisBlock, conn, pack, h.newChain)
	network.Handle(OptAddBlock, conn, pack, h.addBlock)
	// network.Handle(OptGetBlock, conn, pack, getBlock)
	// network.Handle(OptLastHash, conn, pack, getLastHash)
	// network.Handle(OptGetBalance, conn, pack, getBalance)
	// network.Handle(OptGetSizeChain, conn, pack, getChainSize)
}

func (h *Handle) Listen(addr string) {
	network.Listen(addr, h.handleServer)
}

func (h *Handle) newChain(pkg *network.Package) string {
	genesisstr := pkg.Data
	block, err := blockchain.DeserializeBlock(genesisstr)
	if err != nil {
		return err.Error()
	}
	_, err = blockchain.NewBlockchainWithGenesis(block, dbname+h.node.addr, block.Miner)
	if err != nil {
		return err.Error()
	}
	hub.InsertIntoHub(pkg, h.node.hub)
	return Success
}

func (h *Handle) addBlock(pkg *network.Package) string {
	blockstr := pkg.Data
	block, err := blockchain.DeserializeBlock(blockstr)
	if err != nil {
		return err.Error()
	}
	bc, err := blockchain.Load(dbname + h.node.addr)
	if err != nil {
		return err.Error()
	}
	if !block.IsValid(bc) {
		return err.Error()
	}
	err = bc.InsertBlock(block)
	if err != nil {
		return err.Error()
	}
	return Success
}
