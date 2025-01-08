package manage_command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger-transporter/internal/api/grpc/gen/services/ping_pong_gen"
	"slogger-transporter/internal/api/grpc/services/grpc_manager"
	"slogger-transporter/internal/api/grpc/services/ping_pong"
	"slogger-transporter/internal/config"
	"slogger-transporter/pkg/foundation/errs"
	"strings"
	"time"
)

type Command struct {
	managerClient *grpc_manager.Client
	closing       bool
}

func (c *Command) Title() string {
	return "Stop grpc server"
}

func (c *Command) Parameters() string {
	return "{stat || stop}"
}

func (c *Command) Handle(ctx context.Context, arguments []string) error {
	comm, err := c.getCommandByArguments(arguments)

	grpcPort := config.GetConfig().GetGrpcPort()

	c.managerClient, err = grpc_manager.NewClient(":" + grpcPort)

	if err == nil {
		if comm == "stop" {
			err = c.handleStop()
		} else if comm == "stat" {
			err = c.handleStat()
		} else {
			err = errors.New("unknown command " + comm)
		}
	}

	return errs.Err(err)
}

func (c *Command) Close() error {
	var err error

	if c.managerClient != nil {
		err = c.managerClient.Close()
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

func (c *Command) handleStop() error {
	response, err := c.managerClient.Get().Stop(context.TODO(), &gen.GrpcManagerStopRequest{Message: "Stop please"})

	if err != nil {
		return errs.Err(err)
	}

	if !response.Success {
		return errs.Err(errors.New(response.Message))
	}

	_ = c.managerClient.Close()

	pingPongClient, err := ping_pong.NewClient(config.GetConfig().GetGrpcPort())

	if err != nil {
		return errs.Err(err)
	}

	for {
		_, err = pingPongClient.Get().Ping(context.TODO(), &ping_pong_gen.PingPongPingRequest{Message: "ping"})

		if err == nil {
			time.Sleep(1 * time.Second)

			continue
		}

		//if strings.Compare(err.Error(), "rpc error: code = Unavailable desc = error reading from server: EOF") == 0 {
		if strings.Contains(err.Error(), "code = Unavailable desc") {
			slog.Info("Remote application is stopped")

			break
		}

		return errs.Err(err)
	}

	_ = pingPongClient.Close()

	return nil
}

func (c *Command) handleStat() error {
	for {
		if c.closing == true {
			break
		}

		response, err := c.managerClient.Get().Stat(context.TODO(), &gen.GrpcManagerStatRequest{})

		if err != nil {
			return errs.Err(err)
		}

		fmt.Printf("go=%d", response.NumGoroutine)

		fmt.Printf("\tAlloc=%v_MiB", response.AllocMiB)
		fmt.Printf("\tTotalAlloc=%v_MiB", response.TotalAllocMiB)
		fmt.Printf("\tSys=%v_MiB", response.SysMiB)
		fmt.Printf("\tNumGC=%v\n", response.NumGC)

		time.Sleep(1 * time.Second)
	}

	return nil
}
