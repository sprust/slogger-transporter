package queue_listen

import (
	"context"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/queue_service"
	"slogger-transporter/pkg/foundation/errs"
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

func (c *QueueListenCommand) Handle(ctx context.Context, arguments []string) error {
	queueFactory, err := queue_service.NewFactory()

	if err != nil {
		return errs.Err(err)
	}

	for _, queueName := range c.getQueueNames() {
		queue, err := queueFactory.GetQueue(queueName)

		if err != nil {
			return errs.Err(err)
		}

		listener, err := queue_service.NewListener(queue)

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
func (c *QueueListenCommand) getQueueNames() []string {
	return []string{
		config.GetConfig().GetTraceTransporterQueueName(),
	}
}
