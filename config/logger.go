package config

import (
	"github.com/fajarardiyanto/flt-go-logger/interfaces"
	"github.com/fajarardiyanto/flt-go-logger/lib"
	"github.com/fajarardiyanto/go-media-server/internal/model"
)

var logger interfaces.Logger

func init() {
	logger = lib.NewLib()
	logger.Init(model.GetConfig().Name)
}

func GetLogger() interfaces.Logger {
	return logger
}
