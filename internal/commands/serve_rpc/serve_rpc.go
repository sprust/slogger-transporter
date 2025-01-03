package serve_rpc

import (
	"os"
	"slogger-transporter/internal/api/rpc"
	"slogger-transporter/internal/app"
)

type ServeRpcCommand struct {
	server *rpc.Server
}

func (c *ServeRpcCommand) Title() string {
	return "Serve rpc"
}

func (c *ServeRpcCommand) Parameters() string {
	return "{no parameters}"
}

func (c *ServeRpcCommand) Handle(app *app.App, arguments []string) error {
	rpcPort := os.Getenv("RPC_PORT")

	c.server = rpc.NewServer(app, rpcPort)

	err := c.server.Run()

	return err
}

func (c *ServeRpcCommand) Close() error {
	if c.server != nil {
		return c.server.Close()
	}

	return nil
}
