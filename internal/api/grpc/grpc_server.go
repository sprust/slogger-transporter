package grpc

import (
	"google.golang.org/grpc"
	"io"
	"log/slog"
	"net"
	"slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger-transporter/internal/api/grpc/gen/services/ping_pong_gen"
	"slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/api/grpc/services/grpc_manager"
	"slogger-transporter/internal/api/grpc/services/ping_pong"
	"slogger-transporter/internal/api/grpc/services/trace_collector"
	"slogger-transporter/internal/api/grpc/services/trace_transporter"
	"slogger-transporter/internal/app"
)

type Server struct {
	app            *app.App
	rpcPort        string
	sloggerGrpcUrl string
	grpcServer     *grpc.Server
	servers        []io.Closer
}

func NewServer(app *app.App, rpcPort string, sloggerGrpcUrl string) *Server {
	server := &Server{
		app:            app,
		rpcPort:        rpcPort,
		sloggerGrpcUrl: sloggerGrpcUrl,
	}

	app.AddCloseListener(server)

	return server
}

func (s *Server) Run() error {
	s.grpcServer = grpc.NewServer()

	err := s.registerPingPongServer(s.grpcServer)

	if err != nil {
		return err
	}

	err = s.registerTraceCollectorServer(s.grpcServer)

	if err != nil {
		return err
	}

	err = s.registerTraceTransporterServer(s.grpcServer)

	if err != nil {
		return err
	}

	err = s.registerGrpcManagerServer(s.grpcServer)

	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		slog.Error("Error listening: ", err.Error())

		return err
	}

	slog.Info("Listening on " + s.rpcPort)

	err = s.grpcServer.Serve(listener)

	if err != nil {
		slog.Error("Error serving: ", err.Error())
	}

	return err
}

func (s *Server) Close() error {
	if s.grpcServer == nil {
		return nil
	}

	slog.Info("Closing grpc server...")

	for _, server := range s.servers {
		err := server.Close()

		if err != nil {
			return err
		}
	}

	s.grpcServer.Stop()

	return nil
}

func (s *Server) registerPingPongServer(grpcServer *grpc.Server) error {
	server := ping_pong.NewServer()

	ping_pong_gen.RegisterPingPongServer(grpcServer, server)

	s.servers = append(s.servers, server)

	return nil
}

func (s *Server) registerTraceCollectorServer(grpcServer *grpc.Server) error {
	collectorServer, err := trace_collector.NewServer(s.app, s.sloggerGrpcUrl)

	if err != nil {
		slog.Error("Error creating collector: ", err.Error())

		return err
	}

	trace_collector_gen.RegisterTraceCollectorServer(grpcServer, collectorServer)

	s.servers = append(s.servers, collectorServer)

	return nil
}

func (s *Server) registerTraceTransporterServer(grpcServer *grpc.Server) error {
	transporterServer, err := trace_transporter.NewServer(s.app, s.sloggerGrpcUrl)

	if err != nil {
		return err
	}

	trace_transporter_gen.RegisterTraceTransporterServer(grpcServer, transporterServer)

	s.servers = append(s.servers, transporterServer)

	return nil
}

func (s *Server) registerGrpcManagerServer(grpcServer *grpc.Server) error {
	grpcManagerServer := grpc_manager.NewServer(s.app)

	grpc_manager_gen.RegisterGrpcManagerServer(grpcServer, grpcManagerServer)

	s.servers = append(s.servers, grpcManagerServer)

	return nil
}
