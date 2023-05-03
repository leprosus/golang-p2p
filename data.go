package p2p

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Data struct {
	Bytes []byte
	Addr  string
}

func (d *Data) SetBytes(bs []byte) {
	d.Bytes = bs
}

func (d *Data) GetBytes() (bs []byte) {
	return d.Bytes
}

func (d *Data) SetGob(val interface{}) (err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	d.Bytes = buf.Bytes()

	return
}

func (d *Data) GetGob(val interface{}) (err error) {
	err = gob.NewDecoder(bytes.NewReader(d.Bytes)).Decode(val)

	return
}

func (d *Data) SetJson(val interface{}) (err error) {
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	d.Bytes = buf.Bytes()

	return
}

func (d *Data) GetJson(val interface{}) (err error) {
	err = json.NewDecoder(bytes.NewReader(d.Bytes)).Decode(val)

	return
}

func (d *Data) SetAddr(s string) {
	d.Addr = s
}

func (d *Data) String() (str string) {
	return string(d.Bytes)
}
