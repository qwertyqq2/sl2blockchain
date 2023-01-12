package node

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/qwertyqq2/sl2blockchain/crypto"
)

func TestBytesToString(t *testing.T) {
	hash := crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte("adqd"),
			},
			[]byte{},
		))
	encdeHash := crypto.Base64Encode(hash)
	res := crypto.Base64Decode(encdeHash)
	fmt.Println(bytes.Equal(res, hash))
}
