package serve_rpc

import (
	"slogger-transporter/internal/api/rpc"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/errs"
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
	rpcPort := config.GetConfig().GetRpcPort()

	c.server = rpc.NewServer(app, rpcPort)

	err := c.server.Run()

	return errs.Err(err)
}

func (c *ServeRpcCommand) Close() error {
	if c.server != nil {
		return errs.Err(c.server.Close())
	}

	return nil
}
