package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Encode(val interface{}) (bs []byte, err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	bs = buf.Bytes()

	return
}

func MustEncode(val interface{}) (bs []byte) {
	var err error
	bs, err = Encode(val)
	if err != nil {
		log.Panicf("can't encode: %v", err)
	}

	return
}

func Decode(bs []byte) (val interface{}, err error) {
	var buf bytes.Buffer
	err = gob.NewDecoder(&buf).Decode(&val)
	if err != nil {
		return
	}

	_, err = buf.Write(bs)

	return
}

func MustDecode(bs []byte) (val interface{}) {
	var err error
	val, err = Decode(bs)
	if err != nil {
		log.Panicf("can't decode: %v", err)
	}

	return
}
