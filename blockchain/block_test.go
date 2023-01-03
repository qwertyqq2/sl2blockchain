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
