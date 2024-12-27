package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"slogger-transporter/internal/rpc"
	"slogger-transporter/internal/rpc/ping_pong"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// TODO: validate env variables
}

func main() {
	rpcPort := os.Getenv("RPC_PORT")

	functions := []any{
		&ping_pong.PingPong{},
	}

	server := rpc.NewServer(rpcPort, functions)

	err := server.Run()

	if err != nil {
		slog.Error(err.Error())
	}
}
