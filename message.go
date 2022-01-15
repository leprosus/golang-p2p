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

func (msg Message) Encode(ck CipherKey) (cm CryptMessage, err error) {
	var buf bytes.Buffer

	err = gob.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return
	}

	cm, err = ck.Encode(buf.Bytes())

	return
}

func (cm CryptMessage) Decode(ck CipherKey) (msg Message, err error) {
	var bs []byte
	bs, err = ck.Decode(cm)
	if err != nil {
		return
	}

	err = gob.NewDecoder(bytes.NewReader(bs)).Decode(&msg)

	return
}
