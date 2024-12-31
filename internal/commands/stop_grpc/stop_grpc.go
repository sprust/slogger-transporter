package stop_grpc

import (
	"errors"
	"log/slog"
	"os"
	gen "slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger-transporter/internal/api/grpc/services/grpc_manager"
	"slogger-transporter/internal/app"
	"strings"
)

type GrpcStopCommand struct {
}

func (c *GrpcStopCommand) Title() string {
	return "Stop grpc server"
}

func (c *GrpcStopCommand) Parameters() string {
	return "{no parameters}"
}

func (c *GrpcStopCommand) Handle(app *app.App, arguments []string) error {
	grpcPort := os.Getenv("GRPC_PORT")

	client, err := grpc_manager.NewClient(":" + grpcPort)

	if err != nil {
		return err
	}

	_, err = client.Get().Stop(app.GetContext(), &gen.GrpcManagerStopRequest{Message: "Stop please"})

	if err != nil {
		if strings.Compare(err.Error(), "rpc error: code = Unavailable desc = error reading from server: EOF") == 0 {
			slog.Info("Grpc server is stopped")

			return nil
		}

		return err
	}

	// this is never run

	return errors.New("grpc server is not stopped")
}
