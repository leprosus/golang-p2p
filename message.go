package p2p

type Message struct {
	Topic   string
	Content []byte
	Error   error
}
