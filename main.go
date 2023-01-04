package main

import (
	"log"

	"github.com/qwertyqq2/sl2blockchain/blockchain"
)

const (
	DbName = "bc.db"
)

func main() {
	miner := blockchain.NewUser()
	blockchain.NewBlockchain(DbName, miner.Public())
	bc, err := blockchain.Load(DbName)
	if err != nil {
		log.Fatal(err)
	}
	// /bc.PrintBlockchain()
	b := blockchain.NewBlock(miner.Public(), bc.LastBlock().Hash())
	tx1, err := blockchain.NewTransaction(miner, bc.LastBlock().Hash(), "aaa", 3)
	if err != nil {
		log.Fatal(err)
	}
	err = b.InsertTx(bc, tx1)
	if err != nil {
		log.Fatal(err)
	}
	tx2, err := blockchain.NewTransaction(miner, bc.LastBlock().Hash(), "bbb", 3)
	if err != nil {
		log.Fatal(err)
	}
	err = b.InsertTx(bc, tx2)
	if err != nil {
		log.Fatal(err)
	}
	tx3, err := blockchain.NewTransaction(miner, bc.LastBlock().Hash(), "ccc", 3)
	if err != nil {
		log.Fatal(err)
	}
	err = b.InsertTx(bc, tx3)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Accept(bc, miner, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = bc.InsertBlock(b)
	if err != nil {
		log.Fatal(err)
	}
	bc.PrintBlockchain()
}
