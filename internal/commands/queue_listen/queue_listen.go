package queue_listen

import (
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/queue_service"
	"sync"
)

var queueNames = []string{
	queue_service.QueueTraceTransporterName,
}

type QueueListenCommand struct {
	listeners []*queue_service.Listener
}

func (c *QueueListenCommand) Title() string {
	return "Start jobs listening"
}

func (c *QueueListenCommand) Parameters() string {
	return "{no parameters}"
}

func (c *QueueListenCommand) Handle(app *app.App, arguments []string) error {
	queueFactory, err := queue_service.NewFactory(app)

	if err != nil {
		return err
	}

	for _, queueName := range queueNames {
		queue, err := queueFactory.GetQueue(queueName)

		if err != nil {
			return err
		}

		listener, err := queue_service.NewListener(app, queue)

		if err != nil {
			return err
		}

		c.listeners = append(c.listeners, listener)
	}

	waitGroup := sync.WaitGroup{}

	for _, listener := range c.listeners {
		waitGroup.Add(1)

		go func() {
			err := listener.Listen()

			if err != nil {
				panic(err)
			}
		}()
	}

	waitGroup.Wait()

	return nil
}

func (c *QueueListenCommand) Close() error {
	for _, listener := range c.listeners {
		err := listener.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
