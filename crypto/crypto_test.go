package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"
)

func TestRandUint64(t *testing.T) {
	a1 := randuint64()
	a2 := randuint64()
	a3 := randuint64()
	t.Log(a1, a2, a3)
}

func TestLsh(t *testing.T) {
	a := big.NewInt(1)
	b := a.Lsh(a, 256-8)
	t.Log(b)
}

func TestPrintf(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		fmt.Printf("\rval: '%d'", i)
	}
}

func TestSignTx(t *testing.T) {
	pk := GeneratePrivate(2048)
	data := []byte("data: 000101001001001")
	hash := HashSum(
		bytes.Join(
			[][]byte{
				data,
			},
			[]byte{},
		))
	sign, err := Sign(pk, hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Подпись", sign)
	err = Verify(&pk.PublicKey, hash, sign)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("verify")
}

func TestSignWithOpts(t *testing.T) {

	// Generate RSA Keys
	miryanPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		t.Fatal(err)
	}

	miryanPublicKey := &miryanPrivateKey.PublicKey

	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		t.Fatal(err)
	}

	raulPublicKey := &raulPrivateKey.PublicKey

	fmt.Println("Private Key : ", miryanPrivateKey)
	fmt.Println("Public key ", miryanPublicKey)
	fmt.Println("Private Key : ", raulPrivateKey)
	fmt.Println("Public key ", raulPublicKey)

	//Encrypt Miryan Message
	message := []byte("the code must be like a piece of music")
	label := []byte("")
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, raulPublicKey, message, label)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	fmt.Println()

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, miryanPrivateKey, newhash, hashed, &opts)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("PSS Signature : %x\n", signature)

	// Decrypt Message
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, raulPrivateKey, ciphertext, label)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", ciphertext, plainText)

	//Verify Signature
	err = rsa.VerifyPSS(miryanPublicKey, newhash, hashed, signature, &opts)

	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("Verify Signature successful")
	}

}
