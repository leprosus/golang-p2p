package p2p

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

func (pk PublicKey) Encode(bs []byte) (rs []byte, err error) {
	rs, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, &pk.Key, bs, nil)

	return
}

type PrivateKey struct {
	key rsa.PrivateKey
}

func (pk PrivateKey) Decode(bs []byte) (rs []byte, err error) {
	rs, err = pk.key.Decrypt(nil, bs, &rsa.OAEPOptions{Hash: crypto.SHA256})

	return
}
