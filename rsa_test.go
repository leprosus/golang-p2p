package p2p

import "testing"

func TestRSA(t *testing.T) {
	text := []byte("a special secret message")

	rsa, err := NewRSA()
	if err != nil {
		t.Error(err)
	}

	var enc []byte
	enc, err = rsa.PublicKey().Encode(text)
	if err != nil {
		t.Error(err)
	}

	var origin []byte
	origin, err = rsa.PrivateKey().Decode(enc)
	if err != nil {
		t.Error(err)
	}

	if string(text) != string(origin) {
		t.Error("encrypt -> decrypt procedure doesn't work expectedly")
	}
}
