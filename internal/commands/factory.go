package commands

import (
	"errors"
	"slogger-transporter/internal/commands/queue_listen"
	"slogger-transporter/internal/commands/serve_grpc"
	"slogger-transporter/internal/commands/serve_rpc"
	"slogger-transporter/internal/commands/stop_grpc"
)

var commands = map[string]CommandInterface{
	"serve:grpc":   &serve_grpc.ServeGrpcCommand{},
	"grpc:stop":    &stop_grpc.GrpcStopCommand{},
	"serve:rpc":    &serve_rpc.ServeRpcCommand{},
	"queue:listen": &queue_listen.QueueListenCommand{},
}

func GetCommand(name string) (CommandInterface, error) {
	command, ok := commands[name]

	if !ok {
		return nil, errors.New("command not found")
	}

	return command, nil
}

func GetCommands() map[string]CommandInterface {
	return commands
}
