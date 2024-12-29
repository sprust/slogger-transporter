package serve_rpc

import (
	"os"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/temp/rpc"
)

type ServeRpcCommand struct {
}

func (c *ServeRpcCommand) Title() string {
	return "Serve rpc"
}

func (c *ServeRpcCommand) Parameters() string {
	return "{no parameters}"
}

func (c *ServeRpcCommand) Handle(app *app.App, arguments []string) error {
	rpcPort := os.Getenv("RPC_PORT")

	server := rpc.NewServer(rpcPort)

	err := server.Run()

	return err
}
