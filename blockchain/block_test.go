package blockchain

import (
	"testing"
)

func TestBlockEncode(t *testing.T) {
	block := NewGenesisBlock("me")
	serblock, err := SerializeBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(serblock))
}

func TestSerializeBlocks(t *testing.T) {
	block1 := &Block{
		CurHash: GenerateRandom(),
	}
	block2 := &Block{
		CurHash: GenerateRandom(),
	}
	blocks := []*Block{block1, block2}
	blocksStr, err := SerializeBlocks(blocks)
	if err != nil {
		t.Log(err)
	}
	blockcopy, err := DeserializeBlocks(blocksStr)
	if err != nil {
		t.Log(err)
	}
	t.Log(blockcopy)
}
