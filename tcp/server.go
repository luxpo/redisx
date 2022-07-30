package tcp

import (
	"context"
	"net"

	"github.com/luxpo/redisx/interface/tcp"
	"github.com/luxpo/redisx/lib/logger"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	_, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listening...")

	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	ctx := context.Background()
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link")
		go func() {
			handler.Handle(ctx, conn)
		}()
	}
}
