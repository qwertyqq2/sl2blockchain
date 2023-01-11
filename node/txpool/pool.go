package txpool

import (
	"github.com/qwertyqq2/sl2blockchain/blockchain"
)

type Pool struct {
	txs []*blockchain.Transaction
}

func NewPool() *Pool {
	return &Pool{
		txs: make([]*blockchain.Transaction, 0),
	}
}

func (p *Pool) Add(tx *blockchain.Transaction) {
	p.txs = append(p.txs, tx)
}

func (p *Pool) put() *blockchain.Transaction {
	tx := p.txs[0]
	p.txs[0] = nil
	p.txs = p.txs[1:]
	return tx
}

func (p *Pool) GetTxs(maxtx uint) ([]*blockchain.Transaction, bool) {
	txs := make([]*blockchain.Transaction, 0)
	if len(p.txs) > int(maxtx) {
		for i := 0; i < int(maxtx); i++ {
			txs = append(txs, p.put())
		}
		return txs, true
	}
	return nil, false
}
