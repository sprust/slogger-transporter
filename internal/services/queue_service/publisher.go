package queue_service

import "slogger-transporter/internal/app"

type Publisher struct {
	app *app.App
}

func NewPublisher(app *app.App) *Publisher {
	return &Publisher{
		app: app,
	}
}

func (p *Publisher) Publish(queueName string, payload []byte) error {
	queueFactory, err := NewFactory(p.app)

	if err != nil {
		return err
	}

	queue, err := queueFactory.GetQueue(queueName)

	if err != nil {
		return err
	}

	return queue.Publish(payload)
}
