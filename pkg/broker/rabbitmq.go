package broker

import (
	"context"

	"github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/fajarardiyanto/go-media-server/config"
)

func OnMsg(key string, chat interface{}) {
	go func() {
		config.GetRabbitMQ().Producer(interfaces.RabbitMQOptions{NoWait: true})
	}()

	go func() {
		if err := config.GetRabbitMQ().Push(context.Background(), "", key, chat, nil); err != nil {
			config.GetLogger().Error(err)
		}
	}()
}
