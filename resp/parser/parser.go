package parser

import (
	"io"

	"github.com/luxpo/redisx/interface/resp"
)

// Payload stores redis.Reply or error
type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func (this *readState) finished() bool {
	return this.expectedArgsCount > 0 && len(this.args) == this.expectedArgsCount
}

// ParseStream reads data from io.Reader and send payloads through channel
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
}
