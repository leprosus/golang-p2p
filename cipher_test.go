package p2p

import (
	"testing"
)

func TestCipher(t *testing.T) {
	origin := []byte("a special secret message")

	ck, err := NewCipherKey()
	if err != nil {
		t.Fatal(err.Error())
	}

	var encoded []byte
	encoded, err = ck.Encode(origin)
	if err != nil {
		t.Fatal(err.Error())
	}

	var decoded []byte
	decoded, err = ck.Decode(encoded)
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(origin) != string(decoded) {
		t.Fatal("Origin and Decoded are not equal")
	}
}
