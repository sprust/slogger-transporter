package serve_grpc

import (
	"os"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/collector/grpc"
)

type ServeGrpcCommand struct {
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

	server := grpc.NewServer(app, grpcPort, sloggerGrpcUrl)

	err := server.Run()

	return err
}
