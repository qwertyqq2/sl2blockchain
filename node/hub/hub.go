package hub

import (
	"github.com/qwertyqq2/sl2blockchain/blockchain"
	"github.com/qwertyqq2/sl2blockchain/network"
	"github.com/qwertyqq2/sl2blockchain/node/txpool"
)

type Hub struct {
	hub       chan *network.Package
	pool      *txpool.Pool
	addresses []string
	addrnode  string
}

func NewHub(addresses []string, dbname, addr string) *Hub {
	return &Hub{
		addresses: addresses,
		hub:       make(chan *network.Package, len(addresses)),
		pool:      txpool.NewPool(),
		addrnode:  addr,
	}
}

func InsertIntoHub(pkg *network.Package, h *Hub) {
	h.hub <- pkg
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
					}
					//fmt.Println("resp", resp.Option, resp.Data)
				}
			}
		}
	}
}

func (h *Hub) poolCheck(errChan chan error, txsChan chan []*blockchain.Transaction) {
	for {
		txs, f := h.pool.GetTxs()
		if f {
			txsChan <- txs
		}
	}
}

func (h *Hub) ListenHub(errChan chan error, txsChan chan []*blockchain.Transaction) {
	go h.serveNeighbors(errChan)
	go h.poolCheck(errChan, txsChan)
}
