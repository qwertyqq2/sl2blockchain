package blockchain

import (
	"encoding/base64"
	"math/rand"
)

func Base64Encode(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func GenerateRandom() []byte {
	slice := make([]byte, 100)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}
