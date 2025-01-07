package start

import (
	"errors"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/commands/queue_listen"
	"slogger-transporter/internal/commands/serve_grpc"
	"slogger-transporter/internal/services/errs"
	"sync"
)

type Command struct {
	serveGrpcCommand   *serve_grpc.ServeGrpcCommand
	queueListenCommand *queue_listen.QueueListenCommand
}

func (c *Command) Title() string {
	return "Start application"
}

func (c *Command) Parameters() string {
	return "{no parameters}"
}

func (c *Command) Handle(app *app.App, arguments []string) error {
	if len(arguments) != 0 {
		return errs.Err(errors.New("this command does not accept any parameters"))
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		c.serveGrpcCommand = &serve_grpc.ServeGrpcCommand{}

		err := c.serveGrpcCommand.Handle(app, []string{})

		if err != nil {
			panic(errs.Err(err))
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		c.queueListenCommand = &queue_listen.QueueListenCommand{}

		err := c.queueListenCommand.Handle(app, []string{})

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