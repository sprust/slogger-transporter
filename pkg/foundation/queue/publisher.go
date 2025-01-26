package queue

import (
	"slogger/pkg/foundation/config"
	"slogger/pkg/foundation/errs"
	"slogger/pkg/foundation/queue/connections"
	"slogger/pkg/foundation/queue/objects"
)

type Publisher struct {
	config *config.RmqConfig
}

func NewPublisher(config *config.RmqConfig) *Publisher {
	return &Publisher{
		config: config,
	}
}

func (p *Publisher) Publish(settings *objects.QueueSettings, payload []byte) error {
	connection := connections.NewConnection(p.config)

	err := connection.Publish(settings.QueueName, payload)

	_ = connection.Close()

	return errs.Err(err)
}
