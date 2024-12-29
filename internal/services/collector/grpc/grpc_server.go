package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/collector"
	"slogger-transporter/internal/services/collector/grpc/gen/services/ping_pong_gen"
	"slogger-transporter/internal/services/collector/grpc/gen/services/trace_collector_gen"
)

type Server struct {
	app            *app.App
	rpcPort        string
	sloggerGrpcUrl string
}

func NewServer(app *app.App, rpcPort string, sloggerGrpcUrl string) *Server {
	return &Server{
		app:            app,
		rpcPort:        rpcPort,
		sloggerGrpcUrl: sloggerGrpcUrl,
	}
}

func (s *Server) Run() error {
	grpcServer := grpc.NewServer()

	err := s.registerPingPongServer(grpcServer)

	if err != nil {
		return err
	}

	err = s.registerTraceCollectorServer(grpcServer)

	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		slog.Error("Error listening: ", err.Error())

		return err
	}

	slog.Info("Listening on " + s.rpcPort)

	err = grpcServer.Serve(listener)

	if err != nil {
		slog.Error("Error serving: ", err.Error())
	}

	return err
}

func (s *Server) registerPingPongServer(grpcServer *grpc.Server) error {
	ping_pong_gen.RegisterPingPongServer(grpcServer, collector.NewPingPongServer())

	return nil
}

func (s *Server) registerTraceCollectorServer(grpcServer *grpc.Server) error {
	collectorServer, err := collector.NewCollectorServer(s.app, s.sloggerGrpcUrl)

	if err != nil {
		slog.Error("Error creating collector: ", err.Error())

		return err
	}

	trace_collector_gen.RegisterTraceCollectorServer(grpcServer, collectorServer)

	return nil
}
