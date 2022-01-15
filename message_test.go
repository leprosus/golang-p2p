package p2p

import (
	"testing"
)

func TestMessage(t *testing.T) {
	key, err := NewCipherKey()
	if err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Topic:   "topic",
		Content: []byte("some very important text"),
		Error:   nil,
	}

	var cm CryptMessage
	cm, err = msg.Encode(key)
	if err != nil {
		t.Fatal(err)
	}

	var newMsg Message
	newMsg, err = cm.Decode(key)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Topic != newMsg.Topic ||
		string(msg.Content) != string(newMsg.Content) ||
		msg.Error != msg.Error {
		t.Fatal(err)
	}
}
