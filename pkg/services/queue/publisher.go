package queue

import (
	"slogger/pkg/foundation/errs"
	"slogger/pkg/services/queue/connections"
	"slogger/pkg/services/queue/objects"
)

type Publisher struct {
	config objects.RmqConfig
}

func NewPublisher(config objects.RmqConfig) *Publisher {
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
