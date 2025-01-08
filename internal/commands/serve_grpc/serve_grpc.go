package serve_grpc

import (
	"context"
	"slogger-transporter/internal/api/grpc"
	"slogger-transporter/internal/config"
	"slogger-transporter/pkg/foundation/errs"
)

type ServeGrpcCommand struct {
	server *grpc.Server
}

func (c *ServeGrpcCommand) Title() string {
	return "Serve grpc"
}

func (c *ServeGrpcCommand) Parameters() string {
	return "{no parameters}"
}

func (c *ServeGrpcCommand) Handle(ctx context.Context, arguments []string) error {
	cfg := config.GetConfig()

	grpcPort := cfg.GetGrpcPort()
	sloggerGrpcUrl := cfg.GetSloggerGrpcUrl()

	c.server = grpc.NewServer(grpcPort, sloggerGrpcUrl)

	err := c.server.Run()

	return errs.Err(err)
}

func (c *ServeGrpcCommand) Close() error {
	if c.server != nil {
		return errs.Err(c.server.Close())
	}

	return nil
}
