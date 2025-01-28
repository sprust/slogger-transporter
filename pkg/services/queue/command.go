package queue

import (
	"context"
)

type Command struct {
	service *Service
}

func (c *Command) Title() string {
	return "Start jobs listening"
}

func (c *Command) Parameters() string {
	return "{no parameters}"
}

func (c *Command) Handle(ctx context.Context, arguments []string) error {
	c.service = GetQueueService()

	return c.service.Listen()
}

func (c *Command) Close() error {
	if c.service == nil {
		return nil
	}

	return c.service.Close()
}
