package commands

import (
	"slogger/internal/commands/manage_command"
	"slogger/internal/commands/queue_listen_command"
	"slogger/internal/commands/serve_grpc_command"
	"slogger/internal/commands/serve_rpc_command"
	"slogger/internal/commands/start_command"
	foundationCommands "slogger/pkg/foundation/commands"
)

const (
	StartCommandName       = "start"
	ServeGrpcCommandName   = "serve:grpc"
	ManageCommandName      = "manage"
	ServeRpcCommandName    = "serve:rpc"
	QueueListenCommandName = "queue:listen"
)

var commands = map[string]foundationCommands.CommandInterface{
	StartCommandName:       &start_command.Command{},
	ServeGrpcCommandName:   &serve_grpc_command.Command{},
	ManageCommandName:      &manage_command.Command{},
	ServeRpcCommandName:    &serve_rpc_command.Command{},
	QueueListenCommandName: &queue_listen_command.Command{},
}

func GetCommands() map[string]foundationCommands.CommandInterface {
	return commands
}
