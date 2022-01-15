package p2p

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type CipherKey []byte

type CryptCipherKey []byte

func NewCipherKey() (key CipherKey, err error) {
	bs := make([]byte, 10)
	_, err = io.ReadFull(rand.Reader, bs)
	if err != nil {
		return
	}

	sha256Sum := sha256.Sum256(bs)
	md5Sum := md5.Sum(sha256Sum[:])

	key = md5Sum[:]

	return
}

func (key CipherKey) Encode(bs []byte) (rs []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	rs = gcm.Seal(nonce, nonce, bs, nil)

	return
}

func (key CipherKey) Decode(bs []byte) (rs []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := bs[:nonceSize], bs[nonceSize:]

	rs, err = gcm.Open(nil, nonce, cipherText, nil)

	return
}
