package manage

import (
	"errors"
	"fmt"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger-transporter/internal/api/grpc/services/grpc_manager"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
	"strings"
	"time"
)

type Command struct {
	client  *grpc_manager.Client
	closing bool
}

func (c *Command) Title() string {
	return "Stop grpc server"
}

func (c *Command) Parameters() string {
	return "{no parameters}"
}

func (c *Command) Handle(app *app.App, arguments []string) error {
	comm, err := c.getCommandByArguments(arguments)

	grpcPort := app.GetConfig().GetGrpcPort()

	c.client, err = grpc_manager.NewClient(":" + grpcPort)

	if err == nil {
		if comm == "stop" {
			err = c.handleStop(app)
		} else if comm == "stat" {
			err = c.handleStat(app)
		} else {
			err = errors.New("unknown command " + comm)
		}
	}

	return errs.Err(err)
}

func (c *Command) Close() error {
	var err error

	if c.client != nil {
		err = c.client.Close()
	}

	c.closing = true

	return errs.Err(err)
}

func (c *Command) getCommandByArguments(arguments []string) (string, error) {
	if len(arguments) != 1 {
		return "", errs.Err(errors.New("invalid number of arguments"))
	}

	comm := arguments[0]

	if comm == "" {
		return "", errs.Err(errors.New("command is empty"))
	}

	return comm, nil
}

func (c *Command) handleStop(app *app.App) error {
	_, err := c.client.Get().Stop(app.GetContext(), &gen.GrpcManagerStopRequest{Message: "Stop please"})

	if err != nil {
		if strings.Compare(err.Error(), "rpc error: code = Unavailable desc = error reading from server: EOF") == 0 {
			slog.Info("Application is stopped")

			return nil
		}

		return errs.Err(err)
	}

	// this is never run

	return errors.New("grpc server is not stopped")
}

func (c *Command) handleStat(app *app.App) error {
	for {
		if c.closing == true {
			break
		}

		response, err := c.client.Get().Stat(app.GetContext(), &gen.GrpcManagerStatRequest{})

		if err != nil {
			return errs.Err(err)
		}

		fmt.Printf("go: %d", response.NumGoroutine)

		fmt.Printf("\tAlloc = %v MiB", response.AllocMiB)
		fmt.Printf("\tTotalAlloc = %v MiB", response.TotalAllocMiB)
		fmt.Printf("\tSys = %v MiB", response.SysMiB)
		fmt.Printf("\tNumGC = %v\n", response.NumGC)

		time.Sleep(1 * time.Second)
	}

	return nil
}
