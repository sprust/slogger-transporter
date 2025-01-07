package commands

import (
	"errors"
	"slogger-transporter/internal/commands/manage"
	"slogger-transporter/internal/commands/queue_listen"
	"slogger-transporter/internal/commands/serve_grpc"
	"slogger-transporter/internal/commands/serve_rpc"
	"slogger-transporter/internal/services/errs"
)

const (
	ServeGrpcCommandName   = "serve:grpc"
	ManageCommandName      = "manage"
	ServeRpcCommandName    = "serve:rpc"
	QueueListenCommandName = "queue:listen"
)

var commands = map[string]CommandInterface{
	ServeGrpcCommandName:   &serve_grpc.ServeGrpcCommand{},
	ManageCommandName:      &manage.Command{},
	ServeRpcCommandName:    &serve_rpc.ServeRpcCommand{},
	QueueListenCommandName: &queue_listen.QueueListenCommand{},
}

func GetCommand(name string) (CommandInterface, error) {
	command, ok := commands[name]

	if !ok {
		return nil, errs.Err(errors.New("command not found"))
	}

	return command, nil
}

func GetCommands() map[string]CommandInterface {
	return commands
}
