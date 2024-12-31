package serve_grpc

import (
	"os"
	"os/signal"
	"slogger-transporter/internal/api/grpc"
	"slogger-transporter/internal/app"
	"syscall"
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

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

		err := server.Close()

		if err != nil {
			panic(err)
		}
	}()

	err := server.Run()

	return err
}
