package start_command

import (
	"context"
	"errors"
	"slogger-transporter/internal/commands/queue_listen_command"
	"slogger-transporter/internal/commands/serve_grpc_command"
	"slogger-transporter/pkg/foundation/errs"
	"sync"
)

type Command struct {
	serveGrpcCommand   *serve_grpc_command.Command
	queueListenCommand *queue_listen_command.Command
}

func (c *Command) Title() string {
	return "Start application"
}

func (c *Command) Parameters() string {
	return "{no parameters}"
}

func (c *Command) Handle(ctx context.Context, arguments []string) error {
	if len(arguments) != 0 {
		return errs.Err(errors.New("this command does not accept any parameters"))
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		c.serveGrpcCommand = &serve_grpc_command.Command{}

		err := c.serveGrpcCommand.Handle(ctx, []string{})

		if err != nil {
			panic(errs.Err(err))
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		c.queueListenCommand = &queue_listen_command.Command{}

		err := c.queueListenCommand.Handle(ctx, []string{})

		if err != nil {
			panic(errs.Err(err))
		}
	}()

	wg.Wait()

	return nil
}

func (c *Command) Close() error {
	if c.serveGrpcCommand != nil {
		err := c.serveGrpcCommand.Close()

		if err != nil {
			return errs.Err(err)
		}
	}

	if c.queueListenCommand != nil {
		err := c.queueListenCommand.Close()

		if err != nil {
			return errs.Err(err)
		}
	}

	return nil
}
