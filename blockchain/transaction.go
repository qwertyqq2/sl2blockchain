package blockchain

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"

	"github.com/qwertyqq2/sl2blockchain/crypto"
)

const (
	StartPercent = 2
)

type Transaction struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Value     uint64 `json:"value"`
	ToStorage uint64 `json:"toStorage"`
	HashTx    []byte `json:"hashTx"`
	Sing      []byte `json:"sign"`
	RandBytes []byte `json:"randBytes"`
	PrevBlock []byte `json:"prevBlock"`
}

func NewTransaction(user *User, lastHash []byte, to string, value uint64) (*Transaction, error) {
	randBytes := GenerateRandom()
	tx := &Transaction{
		RandBytes: randBytes,
		PrevBlock: lastHash,
		Sender:    user.Public(),
		Receiver:  to,
		Value:     value,
	}
	if value > StartPercent {
		tx.ToStorage = StorageReward
	}
	tx.HashTx = tx.Hash()
	s, err := tx.sign(user.Private())
	if err != nil {
		return nil, err
	}
	tx.Sing = s
	return tx, nil
}

func (tx *Transaction) Hash() []byte {
	return crypto.HashSum(bytes.Join(
		[][]byte{
			tx.RandBytes,
			tx.PrevBlock,
			[]byte(tx.Sender),
			[]byte(tx.Receiver),
			crypto.ToBytes(tx.Value),
			crypto.ToBytes(tx.ToStorage),
		},
		[]byte{},
	))
}

func (tx *Transaction) sign(pk *rsa.PrivateKey) ([]byte, error) {
	return crypto.Sign(pk, tx.HashTx)
}

func (tx *Transaction) IsValid() (bool, error) {
	if !tx.hashIsValid() {
		return false, ErrTxHash
	}
	if !tx.signIsValid() {
		return false, ErrTxSign
	}
	return true, nil
}

func (tx *Transaction) hashIsValid() bool {
	return bytes.Equal(tx.Hash(), tx.HashTx)
}

func (tx *Transaction) signIsValid() bool {
	if err := crypto.Verify(crypto.ParsePublic(tx.Sender), tx.HashTx, tx.Sing); err != nil {
		return false
	}
	return true
}

func SerializeTX(tx *Transaction) (string, error) {
	jsonData, err := json.MarshalIndent(*tx, "", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializeTX(data string) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
