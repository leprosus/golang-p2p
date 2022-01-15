package p2p

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

type RSA struct {
	key *rsa.PrivateKey
}

func NewRSA() (r *RSA, err error) {
	r = &RSA{}

	r.key, err = rsa.GenerateKey(rand.Reader, 2048)

	return
}

func (r *RSA) PublicKey() (pk PublicKey) {
	pk = PublicKey{
		Key: r.key.PublicKey,
	}

	return
}

func (r *RSA) PrivateKey() (pk PrivateKey) {
	pk = PrivateKey{
		key: *r.key,
	}

	return
}

type PublicKey struct {
	Key rsa.PublicKey
}

func (pk PublicKey) Encode(ck CipherKey) (cck CryptCipherKey, err error) {
	cck, err = rsa.EncryptOAEP(sha512.New(), rand.Reader, &pk.Key, ck, nil)

	return
}

type PrivateKey struct {
	key rsa.PrivateKey
}

func (pk PrivateKey) Decode(cck CryptCipherKey) (ck CipherKey, err error) {
	ck, err = rsa.DecryptOAEP(sha512.New(), rand.Reader, &pk.key, cck, nil)

	return
}
