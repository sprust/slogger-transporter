package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"slogger-transporter/internal/grpc/collector"
	"slogger-transporter/internal/grpc/gen/services/ping_pong_gen"
	"slogger-transporter/internal/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/grpc/ping_pong"
)

type Server struct {
	rpcPort        string
	sloggerGrpcUrl string
}

func NewServer(rpcPort string, sloggerGrpcUrl string) *Server {
	return &Server{
		rpcPort:        rpcPort,
		sloggerGrpcUrl: sloggerGrpcUrl,
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

	ping_pong_gen.RegisterPingPongServer(grpcServer, &ping_pong.PingPongServer{})

	collectorServer, err := collector.NewCollector(s.sloggerGrpcUrl)

	if err != nil {
		slog.Error("Error creating collector: ", err.Error())

		return err
	}

	trace_collector_gen.RegisterTraceCollectorServer(grpcServer, collectorServer)

	err = grpcServer.Serve(listener)

	if err != nil {
		slog.Error("Error serving: ", err.Error())
	}

	return err
}
