package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
	gen "slogger-transporter/internal/grpc/gen/services/ping_pong_gen"
	"slogger-transporter/internal/grpc/ping_pong"
)

type Server struct {
	rpcPort   string
	functions []any
}

func NewServer(rpcPort string) *Server {
	return &Server{
		rpcPort: rpcPort,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		slog.Error("Error listening: ", err.Error())

		return err
	}

	slog.Info("Listening on " + s.rpcPort)

	grpcServer := grpc.NewServer()

	gen.RegisterPingPongServer(grpcServer, &ping_pong.PingPongServer{})

	err = grpcServer.Serve(listener)

	if err != nil {
		slog.Error("Error serving: ", err.Error())
	}

	return err
}
