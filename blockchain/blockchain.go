package blockchain

import (
	"bytes"
	"fmt"
)

type Blockchain struct {
	index   uint64
	levelDb *LevelDB
}

func NewBlockchain(filename, receiver string) (*Block, error) {
	l, err := NewLevelDb(filename)
	if err != nil {
		return nil, err
	}
	bc := &Blockchain{
		levelDb: l,
		index:   l.size(),
	}
	genesis := NewGenesisBlock(receiver)
	return genesis, bc.InsertBlock(genesis)
}

func chainExist(filename string) (bool, error) {
	return existLevel(filename)
}

func NewBlockchainWithGenesis(genesis *Block, filename, receiver string) (*Block, error) {
	f, err := chainExist(filename)
	if f {
		return nil, ErrChainAlreadyExist
	}
	l, err := NewLevelDb(filename)
	if err != nil {
		return nil, err
	}
	bc := &Blockchain{
		levelDb: l,
		index:   l.size(),
	}
	return genesis, bc.InsertBlock(genesis)
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
	err = bc.levelDb.insertBlock(block.CurHash, serializeBlock)
	return err
}

func (bc *Blockchain) Balance(address string, size uint64) (uint64, error) {
	return bc.levelDb.balance(address, size)
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.levelDb.lastBlock()
}

func (bc *Blockchain) BlockInChain(block *Block) bool {
	hash := block.CurHash
	blockcopy, err := bc.levelDb.blockByHash(hash)
	if err != nil {
		return false
	}
	if blockcopy == nil {
		return false
	}
	if bytes.Equal(blockcopy.CurHash, block.CurHash) && bytes.Equal(block.PrevHash, blockcopy.PrevHash) {
		return true
	}
	return false
}

func (bc *Blockchain) GetBlocksFromHash(hash []byte) ([]*Block, error) {
	return bc.levelDb.getBlocksFromHash(hash)
}

func (bc *Blockchain) GetBlockAfter(hash []byte) (*Block, error) {
	return bc.levelDb.getBlockAfter(hash)
}

func (bc *Blockchain) PrintBlockchain() {
	blocks, err := bc.levelDb.getBlocks()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Blockchain:")
	for _, b := range blocks {
		fmt.Println("#####Block######")
		fmt.Println(b)
		fmt.Printf("\n\n\n")
	}
}
