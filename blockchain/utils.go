package blockchain

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func Base64Encode(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func GenerateRandom() []byte {
	rand.Seed(time.Now().UnixNano())
	slice := make([]byte, 100)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}
