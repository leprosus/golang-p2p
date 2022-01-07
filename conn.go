package p2p

import (
	"bufio"
)

type Conn struct {
	reader *bufio.Reader
	writer *bufio.Writer
}
