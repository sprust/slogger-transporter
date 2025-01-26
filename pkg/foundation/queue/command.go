package queue

import (
	"context"
)

type Command struct {
}

func (c *Command) Title() string {
	return "Start jobs listening"
}

func (c *Command) Parameters() string {
	return "{no parameters}"
}

func (c *Command) Handle(ctx context.Context, arguments []string) error {
	return GetQueueService().Listen()
}

func (c *Command) Close() error {
	return nil
}
