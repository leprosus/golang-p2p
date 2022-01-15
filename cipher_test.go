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

func BenchmarkNewCipher(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		_, err = NewCipherKey()
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkCipherEncode(b *testing.B) {
	origin := []byte("a special secret message")

	ck, err := NewCipherKey()
	if err != nil {
		b.Fatal(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = ck.Encode(origin)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkCipherDecode(b *testing.B) {
	origin := []byte("a special secret message")

	ck, err := NewCipherKey()
	if err != nil {
		b.Fatal(err.Error())
	}

	var encoded []byte
	encoded, err = ck.Encode(origin)
	if err != nil {
		b.Fatal(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = ck.Decode(encoded)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}
