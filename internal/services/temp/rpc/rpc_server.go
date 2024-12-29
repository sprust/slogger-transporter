package rpc

import (
	"log/slog"
	"net"
	"net/rpc"
	"slogger-transporter/internal/services/temp/rpc/ping_pong"
)

var functions = []any{
	&ping_pong.PingPong{},
}

type Server struct {
	rpcPort string
}

func NewServer(rpcPort string) *Server {
	return &Server{
		rpcPort: rpcPort,
	}
}

func (srv *Server) Run() error {
	for _, function := range functions {
		err := rpc.Register(function)

		if err != nil {
			slog.Error(err.Error())

			return err
		}
	}

	listener, err := net.Listen("tcp", ":"+srv.rpcPort)

	if err != nil {
		slog.Error("Error listening:", err.Error())

		return err
	}

	defer func(listener net.Listener) {
		err := listener.Close()

		if err == nil {
			return
		}

		panic(err)
	}(listener)

	slog.Info("Listening on port " + srv.rpcPort)

	rpc.Accept(listener)

	return nil
}
