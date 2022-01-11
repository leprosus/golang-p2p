package p2p

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Response struct {
	bs []byte
}

func (res *Response) SetBytes(bs []byte) {
	res.bs = bs
}

func (res *Response) GetBytes() (bs []byte) {
	return res.bs
}

func (res *Response) SetGob(val interface{}) (err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	res.bs = buf.Bytes()

	return
}

func (res *Response) GetGob(val interface{}) (err error) {
	err = gob.NewDecoder(bytes.NewReader(res.bs)).Decode(val)

	return
}

func (res *Response) SetJson(val interface{}) (err error) {
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	res.bs = buf.Bytes()

	return
}

func (res *Response) GetJson(val interface{}) (err error) {
	err = json.NewDecoder(bytes.NewReader(res.bs)).Decode(val)

	return
}

func (res *Response) String() (str string) {
	return string(res.bs)
}
