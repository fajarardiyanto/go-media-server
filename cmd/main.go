package main

import (
	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		config.GetLogger().Error(err).Quit()
	}
}
