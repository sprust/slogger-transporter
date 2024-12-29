package commands

import (
	"errors"
	"slogger-transporter/internal/commands/serve_grpc"
	"slogger-transporter/internal/commands/serve_rpc"
)

var commands = map[string]CommandInterface{
	"serve:grpc": &serve_grpc.ServeGrpcCommand{},
	"serve:rpc":  &serve_rpc.ServeRpcCommand{},
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
