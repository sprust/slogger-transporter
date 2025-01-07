package queue_listen

import (
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
	"slogger-transporter/internal/services/queue_service"
	"sync"
)

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
		return errs.Err(err)
	}

	for _, queueName := range c.getQueueNames(app) {
		queue, err := queueFactory.GetQueue(queueName)

		if err != nil {
			return errs.Err(err)
		}

		listener, err := queue_service.NewListener(app, queue)

		if err != nil {
			return errs.Err(err)
		}

		c.listeners = append(c.listeners, listener)
	}

	waitGroup := sync.WaitGroup{}

	for _, listener := range c.listeners {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

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
			return errs.Err(err)
		}
	}

	return nil
}
func (c *QueueListenCommand) getQueueNames(app *app.App) []string {
	config := app.GetConfig()

	return []string{
		config.GetTraceTransporterQueueName(),
	}
}
