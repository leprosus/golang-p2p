package p2p

import "testing"

func TestRSA(t *testing.T) {
	origin := []byte("a special secret message")

	rsa, err := NewRSA()
	if err != nil {
		t.Fatal(err)
	}

	var encoded []byte
	encoded, err = rsa.PublicKey().Encode(origin)
	if err != nil {
		t.Fatal(err)
	}

	var decoded []byte
	decoded, err = rsa.PrivateKey().Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if string(origin) != string(decoded) {
		t.Fatal("Origin and Decoded are not equal")
	}
}

func BenchmarkNewRSA(b *testing.B) {
	var err error

	_, err = NewRSA()
	for i := 0; i < b.N; i++ {
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRSAEncode(b *testing.B) {
	origin := []byte("a special secret message")

	rsa, err := NewRSA()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = rsa.PublicKey().Encode(origin)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRSADecode(b *testing.B) {
	origin := []byte("a special secret message")

	rsa, err := NewRSA()
	if err != nil {
		b.Fatal(err)
	}

	var encoded []byte
	encoded, err = rsa.PublicKey().Encode(origin)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = rsa.PrivateKey().Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
