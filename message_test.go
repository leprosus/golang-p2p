package p2p

import "testing"

func TestMessage(t *testing.T) {
	rsa, err := NewRSA()
	if err != nil {
		t.Error(err)
	}

	msg := Message{
		Topic:   "topic",
		Content: []byte("some very important text"),
		Error:   nil,
	}

	var cm CryptMessage
	cm, err = msg.Encode(rsa.PublicKey())
	if err != nil {
		t.Error(err)
	}

	var newMsg Message
	newMsg, err = cm.Decode(rsa.PrivateKey())
	if err != nil {
		t.Error(err)
	}

	if msg.Topic != newMsg.Topic ||
		string(msg.Content) != string(newMsg.Content) ||
		msg.Error != msg.Error {
		t.Error(err)
	}
}
