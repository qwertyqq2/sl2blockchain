package hub

import (
	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
	"github.com/qwertyqq2/sl2blockchain/node/txpool"
)

type Hub struct {
	hub         chan *network.Package
	pool        *txpool.Pool
	addresses   []string
	addrnode    string
	PendingResp chan *network.Package
}

func NewHub(addresses []string, dbname, addr string) *Hub {
	return &Hub{
		addresses:   addresses,
		hub:         make(chan *network.Package, len(addresses)),
		pool:        txpool.NewPool(),
		addrnode:    addr,
		PendingResp: make(chan *network.Package, 5),
	}
}

func InsertIntoHub(pkg *network.Package, h *Hub) {
	h.hub <- pkg
}

func InsertIntoPool(tx *blockchain.Transaction, h *Hub) {
	h.pool.Add(tx)
}

func (h *Hub) serveNeighbors(errChan chan error) {
	for {
		select {
		case pack := <-h.hub:
			for _, addr := range h.addresses {
				if addr != h.addrnode {
					resp := network.Send("0.0.0.0"+addr, pack)
					if resp == nil {
						errChan <- ErrNilPackageResp
					} else {
						h.PendingResp <- resp
					}
				}
			}
		}
	}
}

func (h *Hub) poolCheck(errChan chan error, txsChan chan []*blockchain.Transaction, maxtx uint) {
	for {
		txs, f := h.pool.GetTxs(maxtx)
		if f {
			txsChan <- txs
		}
	}
}

func (h *Hub) ListenHub(errChan chan error, txsChan chan []*blockchain.Transaction, maxtx uint) {
	go h.serveNeighbors(errChan)
	go h.poolCheck(errChan, txsChan, maxtx)
}
