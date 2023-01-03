package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

func HashSum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func Sign(pk *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, pk, crypto.SHA256, data, nil)
}

func Verify(pub *rsa.PublicKey, data, sign []byte) error {
	return rsa.VerifyPSS(pub, crypto.SHA256, data, sign, nil)
}

func ParsePrivate(privData string) *rsa.PrivateKey {
	pub, err := x509.ParsePKCS1PrivateKey(Base64Decode(privData))
	if err != nil {
		return nil
	}
	return pub
}

func Base64Decode(data string) []byte {
	result, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return result
}

func ProowOfWork(blockHash []byte, diff uint8, ch chan bool) (uint64, bool) {
	var (
		Target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = randuint64()
		hash    []byte
	)
	Target.Lsh(Target, 256-uint(diff))
	fmt.Println("Start mining")
	for nonce < math.MaxUint64 {
		select {
		case <-ch:
			return 0, false
		default:
			hash = HashSum(bytes.Join(
				[][]byte{
					blockHash,
					ToBytes(nonce),
				},
				[]byte{},
			))
			intHash.SetBytes(hash)
			fmt.Printf("\rproof: %d", intHash)
			if intHash.Cmp(Target) == -1 {
				return nonce, true
			}
			nonce += 1
		}
	}
	return nonce, true
}

func ToBytes(num uint64) []byte {
	data := new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}
	return data.Bytes()
}

func randuint64() uint64 {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(b)
}

func GeneratePrivate(bits uint) *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil
	}
	return priv
}

func StringPublic(pub *rsa.PublicKey) string {
	return Base64Encode(x509.MarshalPKCS1PublicKey(pub))
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ParsePublic(pubData string) *rsa.PublicKey {
	pub, err := x509.ParsePKCS1PublicKey(Base64Decode(pubData))
	if err != nil {
		return nil
	}
	return pub
}
