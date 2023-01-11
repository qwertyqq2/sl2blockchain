package blockchain

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/qwertyqq2/sl2blockchain/crypto"
)

type Block struct {
	CurHash      []byte            `json:"curHash"`
	PrevHash     []byte            `json:"prevHash"`
	Nonce        uint64            `json:"nonce"`
	Difficulty   uint8             `json:"diff"`
	Miner        string            `json:"miner"`
	Sign         []byte            `json:"sign"`
	Timestamp    string            `json:"timestamp"`
	Transactions []Transaction     `json:"transactions"`
	Mapping      map[string]uint64 `json:"mapping"`
}

const (
	GenesisBlock  = "GENESISBLOCK"
	StorageValue  = 100
	GenesisRevard = 100
	StorageChain  = "StorageChain"
	Difficulty    = 10
	StorageReward = 1
	Separator     = "__A__"
)

func SerializeBlock(block *Block) (string, error) {
	jsonData, err := json.MarshalIndent(block, " ", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializeBlock(data string) (*Block, error) {
	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func SerializeBlocks(blocks []*Block) (string, error) {
	res := ""
	for _, b := range blocks {
		bstr, err := SerializeBlock(b)
		if err != nil {
			return "", err
		}
		res += bstr + Separator
	}
	return res, nil
}

func DeserializeBlocks(blocksstr string) ([]*Block, error) {
	blocksSplit := strings.Split(blocksstr, Separator)
	blocks := make([]*Block, 0)
	for _, bstr := range blocksSplit {
		b, err := DeserializeBlock(bstr)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func NewGenesisBlock(receiver string) *Block {
	genesis := &Block{
		PrevHash:  []byte(GenesisBlock),
		Mapping:   make(map[string]uint64),
		Miner:     receiver,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	genesis.Mapping[StorageChain] = StorageValue
	genesis.Mapping[receiver] = GenesisRevard
	genesis.CurHash = genesis.Hash()
	return genesis
}

func NewBlock(miner string, prevhash []byte) *Block {
	return &Block{
		Difficulty:   Difficulty,
		PrevHash:     prevhash,
		Miner:        miner,
		Mapping:      make(map[string]uint64),
		Transactions: make([]Transaction, 0),
	}
}

func (block *Block) Accept(bc *Blockchain, u *User, ch chan bool) error {
	toMiner := uint64(0)
	for _, tx := range block.Transactions {
		toMiner += tx.ToStorage
	}
	block.InsertTx(bc, &Transaction{
		RandBytes: GenerateRandom(),
		PrevBlock: bc.LastBlock().CurHash,
		Sender:    StorageChain,
		Receiver:  u.Public(),
		Value:     toMiner,
	})
	if f, err := block.validTransactions(bc, bc.Size()); !f {
		return err
	}
	block.Timestamp = time.Now().Format(time.RFC3339)
	block.CurHash = block.Hash()
	block.Sign = block.signature(u.Private())
	n, err := block.proof(ch)
	if err != nil {
		return err
	}
	block.Nonce = n
	return nil

}

func (block *Block) Hash() []byte {
	temp := []byte{}
	for _, tx := range block.Transactions {
		temp = crypto.HashSum(
			bytes.Join(
				[][]byte{
					temp,
					tx.HashTx,
				},
				[]byte{},
			))
	}
	list := []string{}
	for addr := range block.Mapping {
		list = append(list, addr)
	}
	sort.Strings(list)
	for _, addr := range list {
		temp = crypto.HashSum(
			bytes.Join(
				[][]byte{
					temp,
					[]byte(addr),
					crypto.ToBytes(block.Mapping[addr]),
				},
				[]byte{},
			))
	}

	return crypto.HashSum(
		bytes.Join(
			[][]byte{
				temp,
				crypto.ToBytes(uint64(block.Difficulty)),
				block.PrevHash,
				[]byte(block.Miner),
				[]byte(block.Timestamp),
			},
			[]byte{},
		))
}

func (block *Block) proof(ch chan bool) (uint64, error) {
	nonce, f := crypto.ProowOfWork(block.CurHash, block.Difficulty, ch)
	if !f {
		return 0, ErrNotProof
	}
	return nonce, nil

}

func (block *Block) signature(pk *rsa.PrivateKey) []byte {
	s, err := crypto.Sign(pk, block.CurHash)
	if err != nil {
		log.Println(err)
		return nil
	}
	return s
}

func (block *Block) validHash() bool {
	return bytes.Equal(block.Hash(), block.CurHash)
}

func (block *Block) validId(bc *Blockchain, id uint64) bool {
	idscan, err := bc.levelDb.idByHash(Base64Encode(block.PrevHash))
	if err != nil {
		log.Println(err)
		return false
	}
	return idscan == id
}

// Валидность т-ции проверяется после вставки в блок
func (block *Block) InsertTx(bc *Blockchain, tx *Transaction) error {
	if tx == nil {
		return ErrNilTx
	}
	var balanceInChain uint64
	balanceTx := tx.Value + tx.ToStorage
	if val, ok := block.Mapping[tx.Sender]; ok {
		balanceInChain = val
	} else {
		bal, err := bc.Balance(tx.Sender, bc.Size())
		if err != nil {
			return err
		}
		balanceInChain = bal
	}
	if balanceInChain < balanceTx {
		return ErrNotEnoghtMoney
	}
	block.Mapping[tx.Sender] = balanceInChain - balanceTx
	block.addBalance(bc, tx.Receiver, tx.Value)
	block.addBalance(bc, StorageChain, tx.ToStorage)
	block.Transactions = append(block.Transactions, *tx)
	return nil
}

func (block *Block) addBalance(bc *Blockchain, receiver string, value uint64) error {
	var balanceInChain uint64
	if val, ok := block.Mapping[receiver]; ok {
		balanceInChain = val
	} else {
		bal, err := bc.Balance(receiver, bc.Size())
		if err != nil {
			return err
		}
		balanceInChain = bal
	}
	block.Mapping[receiver] = balanceInChain + value
	return nil
}

// ErrIncorrectBalanceBlock
func (block *Block) validTransactions(bc *Blockchain, size uint64) (bool, error) {
	lenBlock := len(block.Transactions)
	inStorage := false
	if lenBlock == 0 {
		return false, ErrNilBlock
	}
	for _, tx := range block.Transactions {
		if tx.Sender == StorageChain {
			inStorage = true
			break
		}
	}
	if !inStorage {
		return false, ErrNothaveStorage
	}
	for i := 0; i < lenBlock-1; i++ {
		for j := i + 1; j < lenBlock; j++ {
			if bytes.Equal(block.Transactions[i].RandBytes, block.Transactions[j].RandBytes) {
				return false, ErrEqualRandBytes
			}
			if block.Transactions[i].Sender == StorageChain &&
				block.Transactions[j].Sender == StorageChain {
				return false, ErrSecondStorageSender
			}
		}
	}
	toMiner := uint64(0)
	for _, tx := range block.Transactions {
		toMiner += tx.ToStorage
	}
	for _, tx := range block.Transactions {
		if tx.Sender == StorageChain {
			if tx.Receiver != block.Miner || tx.Value != toMiner {
				return false, ErrIncorrectStorageReceiver
			}
		} else {
			if f, err := tx.IsValid(); !f {
				return f, err
			}
			if f, err := block.validBalance(bc, tx.Sender, bc.Size()); !f {
				return f, err
			}
			if f, err := block.validBalance(bc, tx.Receiver, bc.Size()); !f {
				return f, err
			}

		}
	}
	return true, nil
}

func (block *Block) validTx(bc *Blockchain, size uint64) bool {
	f, err := block.validTransactions(bc, size)
	if !f {
		log.Println(err)
		return false
	}
	return f
}

func (block *Block) validBalance(bc *Blockchain, address string, size uint64) (bool, error) {
	if _, ok := block.Mapping[address]; !ok {
		return false, ErrMissingAddressInBlock
	}
	balanceInChain, err := bc.Balance(address, size)
	sub := uint64(0)
	add := uint64(0)
	if err != nil {
		return false, err
	}
	for _, tx := range block.Transactions {
		if tx.Sender == address {
			sub += tx.Value + tx.ToStorage
		}
		if tx.Receiver == address {
			add += tx.Value
		}
		if address == StorageChain {
			add += tx.ToStorage
		}
	}
	resBal := balanceInChain - sub + add
	if resBal != block.Mapping[address] {
		return false, ErrIncorrectBalanceBlock
	}
	return true, nil
}

func (block *Block) validSign() bool {
	if err := crypto.Verify(crypto.ParsePublic(block.Miner), block.CurHash, block.Sign); err != nil {
		return false
	}
	return true
}

func (block *Block) validProof() bool {
	Target := big.NewInt(1)
	intHash := big.NewInt(1)
	hash := crypto.HashSum(
		bytes.Join(
			[][]byte{
				block.CurHash,
				crypto.ToBytes(block.Nonce),
			},
			[]byte{},
		))
	Target.Lsh(Target, 256-uint(block.Difficulty)) //поч uint
	intHash.SetBytes(hash)
	if intHash.Cmp(Target) == -1 {
		return true
	}
	return false
}

func (block *Block) validMapping() bool {
	for addr := range block.Mapping {
		if addr == StorageChain {
			continue
		}
		f := false
		for _, tx := range block.Transactions {
			if tx.Sender == addr || tx.Receiver == addr {
				f = true
				break
			}
		}
		if !f {
			return false
		}
	}
	return true
}

func (block *Block) validTimestamp(bc *Blockchain) (bool, error) {
	t, err := time.Parse(time.RFC3339, block.Timestamp)
	if err != nil {
		return false, err
	}
	b, err := bc.levelDb.blockByHash(Base64Encode(block.PrevHash))
	if b == nil {
		return false, ErrNilBlock
	}
	tb, err := time.Parse(time.RFC3339, b.Timestamp)
	if err != nil {
		return false, err
	}
	if t.Sub(tb) < 0 {
		return false, ErrIncorrectTimeBlock
	}
	return true, nil
}

func (block *Block) validTime(bc *Blockchain) bool {
	f, err := block.validTimestamp(bc)
	if !f {
		log.Println(err)
	}
	return f
}

func (block *Block) IsValid(bc *Blockchain) error {
	switch {
	case block == nil:
		return ErrNilBlock
	case block.Difficulty != Difficulty:
		return ErrIncorrectDiff
	case !block.validHash():
		return ErrIncorrectHash
	case !block.validTx(bc, bc.Size()):
		return ErrIncorrectTrasnactions
	case !block.validSign():
		return ErrIncorrectSign
	case !block.validProof():
		return ErrIncorrectProof
	case !block.validTime(bc):
		return ErrIncorrectTime
	}
	return nil
}

func (block *Block) Print() {
	fmt.Println(string(block.CurHash), string(block.PrevHash), block.Miner, block.Mapping, block.Nonce, block.Miner, block.Timestamp)
}
