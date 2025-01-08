package queue_service

import (
	"slogger-transporter/pkg/foundation/errs"
)

type Publisher struct {
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(queueName string, payload []byte) error {
	queueFactory, err := NewFactory()

	if err != nil {
		return errs.Err(err)
	}

	queue, err := queueFactory.GetQueue(queueName)

	if err != nil {
		return errs.Err(err)
	}

	return queue.Publish(payload)
}
