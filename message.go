package p2p

import (
	"bytes"
	"encoding/gob"
)

type Message struct {
	Topic   string
	Content []byte
	Error   error
}

type CryptMessage []byte

func (msg Message) Encode(pk PublicKey) (cm CryptMessage, err error) {
	var buf bytes.Buffer

	err = gob.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return
	}

	cm, err = pk.Encode(buf.Bytes())

	return
}

func (cm CryptMessage) Decode(pk PrivateKey) (msg Message, err error) {
	var bs []byte
	bs, err = pk.Decode(cm)
	if err != nil {
		return
	}

	err = gob.NewDecoder(bytes.NewReader(bs)).Decode(&msg)

	return
}
