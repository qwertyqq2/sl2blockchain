package blockchain

import (
	"crypto/rsa"

	"github.com/qwertyqq2/sl2blockchain/crypto"
)

const (
	KeySize = 2048
)

type User struct {
	privateKey *rsa.PrivateKey
}

func NewUser() *User {
	return &User{
		privateKey: crypto.GeneratePrivate(KeySize),
	}
}

func ParseUser(pk string) *User {
	pkkey := crypto.ParsePrivate(pk)
	if pkkey == nil {
		return nil
	}
	return &User{
		privateKey: pkkey,
	}

}

func (user *User) Public() string {
	return crypto.StringPublic(user.PrivateToPublic())
}

func (user *User) Private() *rsa.PrivateKey {
	return user.privateKey
}

func (user *User) PrivateToPublic() *rsa.PublicKey {
	pub := user.privateKey.PublicKey
	return &pub
}
