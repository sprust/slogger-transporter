package rpc

import (
	"log/slog"
	"net"
	"net/rpc"
)

type Server struct {
	rpcPort   string
	functions []any
}

func NewServer(rpcPort string, functions []any) *Server {
	return &Server{
		rpcPort:   rpcPort,
		functions: functions,
	}
}

func (srv *Server) Run() error {
	for _, function := range srv.functions {
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
