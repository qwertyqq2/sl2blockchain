package blockchain

import (
	"fmt"
)

type Blockchain struct {
	index   uint64
	levelDb *LevelDB
}

func NewBlockchain(filename, receiver string) error {
	l, err := NewLevelDb(filename)
	if err != nil {
		return err
	}
	bc := &Blockchain{
		levelDb: l,
		index:   l.size(),
	}
	genesis := NewGenesisBlock(receiver)
	return bc.InsertBlock(genesis)
}

func Load(filename string) (*Blockchain, error) {
	bc, err := loadBlockchain(filename)
	if err != nil {
		return nil, err
	}
	return bc, nil
}

func (bc *Blockchain) Size() uint64 {
	return bc.levelDb.size()
}

func (bc *Blockchain) InsertBlock(block *Block) error {
	bc.index += 1
	serializeBlock, err := SerializeBlock(block)
	if err != nil {
		return err
	}
	err = bc.levelDb.insertBlock(Base64Encode(block.CurHash), serializeBlock)
	return err
}

func (bc *Blockchain) Balance(address string, size uint64) (uint64, error) {
	return bc.levelDb.balance(address, size)
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.levelDb.lastBlock()
}

func (bc *Blockchain) PrintBlockchain() {
	blocks, err := bc.levelDb.getBlocks()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Blockchain:")
	for _, b := range blocks {
		b.Print()
	}
}
