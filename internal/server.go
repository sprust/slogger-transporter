package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"slogger-transporter/internal/grpc"
	"slogger-transporter/internal/rpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// TODO: validate env variables
}

func main() {
	args := os.Args
	argsLen := len(args)

	errorArgMsg := "use 'rpc' or 'grpc' arg"

	if argsLen != 2 {
		slog.Error(errorArgMsg)

		return
	}

	mode := args[1]

	if mode != "rpc" && mode != "grpc" {
		slog.Error(errorArgMsg)

		return
	}

	slog.Warn("*** " + mode + " ***")

	if mode == "rpc" {
		rpcPort := os.Getenv("RPC_PORT")

		server := rpc.NewServer(rpcPort)

		err := server.Run()

		if err != nil {
			slog.Error(err.Error())
		}
	} else {
		grpcPort := os.Getenv("GRPC_PORT")

		server := grpc.NewServer(grpcPort)

		err := server.Run()

		if err != nil {
			slog.Error(err.Error())
		}
	}

	slog.Info("exit")
}
