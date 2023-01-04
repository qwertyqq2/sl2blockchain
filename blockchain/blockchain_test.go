package blockchain

import (
	"log"
	"testing"
)

const (
	DbName = "bctest.db"
)

func TestInsertBlock(t *testing.T) {
	miner := NewUser()
	NewBlockchain(DbName, miner.Public())
	bc, err := Load(DbName)
	if err != nil {
		log.Fatal(err)
	}
	// /bc.PrintBlockchain()
	b := NewBlock(miner.Public(), bc.LastBlock().Hash())
	tx1, err := NewTransaction(miner, bc.LastBlock().Hash(), "me1", 3)
	if err != nil {
		log.Fatal(err)
	}
	err = b.InsertTx(bc, tx1)
	if err != nil {
		log.Fatal(err)
	}
	tx2, err := NewTransaction(miner, bc.LastBlock().Hash(), "me2", 3)
	if err != nil {
		log.Fatal(err)
	}
	err = b.InsertTx(bc, tx2)
	if err != nil {
		log.Fatal(err)
	}
	tx3, err := NewTransaction(miner, bc.LastBlock().Hash(), "me3", 3)
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
