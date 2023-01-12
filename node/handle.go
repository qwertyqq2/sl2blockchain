package node

import (
	"fmt"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/crypto"
	"github.com/qwertyqq2/sl2blockchain/network"
	"github.com/qwertyqq2/sl2blockchain/node/hub"
)

var (
	OptGenesisBlock = 200
	OptAddBlock     = 201
	OptGetBlocks    = 202
	OptLastHash     = 203
	OptGetBalance   = 204
	OptGetSizeChain = 205
	OptNewTx        = 206
	Success         = "success"
	Wait            = "wait"
	Pending         = "pending"
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
	network.Handle(OptNewTx, conn, pack, h.newTx)
	network.Handle(OptGetBlocks, conn, pack, h.getBlocksFromHash)
	//network.Handle(OptGetBlock, conn, pack, h.getBlock)
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
	err = h.node.AddBlock(block, bc)
	if err != nil {
		return err.Error()
	}
	return Success
}

func (h *Handle) newTx(pkg *network.Package) string {
	txstr := pkg.Data
	tx, err := blockchain.DeserializeTX(txstr)
	if err != nil {
		return err.Error()
	}
	hub.InsertIntoPool(tx, h.node.hub)
	return Pending
}

func (h *Handle) getBlocksFromHash(pkg *network.Package) string {
	hash := crypto.Base64Decode(pkg.Data)
	bc, err := blockchain.Load(dbname + h.node.addr)
	if err != nil {
		return err.Error()
	}
	blocks, err := bc.GetBlocksFromHash(hash)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	blocksStr, err := blockchain.SerializeBlocks(blocks)
	if err != nil {
		return err.Error()
	}
	return blocksStr

}
