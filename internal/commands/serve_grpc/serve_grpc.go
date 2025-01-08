package serve_grpc

import (
	"slogger-transporter/internal/api/grpc"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/errs"
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

func (c *ServeGrpcCommand) Handle(app *app.App, arguments []string) error {
	cfg := config.GetConfig()

	grpcPort := cfg.GetGrpcPort()
	sloggerGrpcUrl := cfg.GetSloggerGrpcUrl()

	c.server = grpc.NewServer(app, grpcPort, sloggerGrpcUrl)

	err := c.server.Run()

	return errs.Err(err)
}

func (c *ServeGrpcCommand) Close() error {
	if c.server != nil {
		return errs.Err(c.server.Close())
	}

	return nil
}
