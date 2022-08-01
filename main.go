package main

import (
	"fmt"
	"os"

	"github.com/luxpo/redisx/config"
	"github.com/luxpo/redisx/lib/logger"
	"github.com/luxpo/redisx/resp/handler"
	"github.com/luxpo/redisx/tcp"
)

const configFile string = "redis.conf"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "redisx",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	}

	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d",
				config.Properties.Bind,
				config.Properties.Port),
		},
		// tcp.NewEchoHandler(),
		handler.MakeHandler(),
	)
	if err != nil {
		logger.Error(err)
	}
}
