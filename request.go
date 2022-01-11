package p2p

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Request struct {
	bs []byte
}

func (req *Request) SetBytes(bs []byte) {
	req.bs = bs
}

func (req *Request) GetBytes() (bs []byte) {
	return req.bs
}

func (req *Request) SetGob(val interface{}) (err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	req.bs = buf.Bytes()

	return
}

func (req *Request) GetGob(val interface{}) (err error) {
	err = gob.NewDecoder(bytes.NewReader(req.bs)).Decode(val)

	return
}

func (req *Request) SetJson(val interface{}) (err error) {
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	req.bs = buf.Bytes()

	return
}

func (req *Request) GetJson(val interface{}) (err error) {
	err = json.NewDecoder(bytes.NewReader(req.bs)).Decode(val)

	return
}

func (req *Request) String() (str string) {
	return string(req.bs)
}
