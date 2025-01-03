package serve_grpc

import (
	"os"
	"slogger-transporter/internal/api/grpc"
	"slogger-transporter/internal/app"
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
	grpcPort := os.Getenv("GRPC_PORT")
	sloggerGrpcUrl := os.Getenv("SLOGGER_SERVER_GRPC_URL")

	c.server = grpc.NewServer(app, grpcPort, sloggerGrpcUrl)

	err := c.server.Run()

	return err
}

func (c *ServeGrpcCommand) Close() error {
	if c.server != nil {
		return c.server.Close()
	}

	return nil
}
