package handler

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/luxpo/redisx/database"
	idatabase "github.com/luxpo/redisx/interface/database"
	"github.com/luxpo/redisx/lib/logger"
	"github.com/luxpo/redisx/lib/sync/atomic"
	"github.com/luxpo/redisx/resp/connection"
	"github.com/luxpo/redisx/resp/parser"
	"github.com/luxpo/redisx/resp/reply"
)

/*
 * A tcp.RespHandler implements redis protocol
 */

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// RespHandler implements tcp.Handler and serves as a redis handler
type RespHandler struct {
	activeConn sync.Map // *client -> placeholder
	db         idatabase.Database
	closing    atomic.Boolean // refusing new client and new request
}

// MakeHandler creates a RespHandler instance
func MakeHandler() *RespHandler {
	var db idatabase.Database
	db = database.NewEchoDatabase()
	return &RespHandler{
		db: db,
	}
}

func (h *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}

// Handle receives and executes redis commands
func (h *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
	}

	client := connection.NewConn(conn)
	h.activeConn.Store(client, 1)

	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// connection closed
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// protocol err
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := h.db.Exec(client, r.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unknownErrReplyBytes)
		}
	}
}

// Close stops handler
func (h *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	// TODO: concurrent wait
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}
