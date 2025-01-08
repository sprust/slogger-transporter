package commands

import (
	"slogger-transporter/internal/commands/manage"
	"slogger-transporter/internal/commands/queue_listen"
	"slogger-transporter/internal/commands/serve_grpc"
	"slogger-transporter/internal/commands/serve_rpc"
	"slogger-transporter/internal/commands/start"
	foundationCommands "slogger-transporter/pkg/foundation/commands"
)

const (
	StartCommandName       = "start"
	ServeGrpcCommandName   = "serve:grpc"
	ManageCommandName      = "manage"
	ServeRpcCommandName    = "serve:rpc"
	QueueListenCommandName = "queue:listen"
)

var commands = map[string]foundationCommands.CommandInterface{
	StartCommandName:       &start.Command{},
	ServeGrpcCommandName:   &serve_grpc.ServeGrpcCommand{},
	ManageCommandName:      &manage.Command{},
	ServeRpcCommandName:    &serve_rpc.ServeRpcCommand{},
	QueueListenCommandName: &queue_listen.QueueListenCommand{},
}

func GetCommands() map[string]foundationCommands.CommandInterface {
	return commands
}
