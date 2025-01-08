package serve_grpc_command

import (
	"context"
	"slogger-transporter/internal/api/grpc"
	"slogger-transporter/internal/config"
	"slogger-transporter/pkg/foundation/errs"
)

type Command struct {
	server *grpc.Server
}

func (c *Command) Title() string {
	return "Serve grpc"
}

func (c *Command) Parameters() string {
	return "{no parameters}"
}

func (c *Command) Handle(ctx context.Context, arguments []string) error {
	cfg := config.GetConfig()

	grpcPort := cfg.GetGrpcPort()
	sloggerGrpcUrl := cfg.GetSloggerGrpcUrl()

	c.server = grpc.NewServer(grpcPort, sloggerGrpcUrl)

	err := c.server.Run()

	return errs.Err(err)
}

func (c *Command) Close() error {
	if c.server != nil {
		return errs.Err(c.server.Close())
	}

	return nil
}
