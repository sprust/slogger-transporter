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
	"slogger-transporter/internal/services/errs"
)

type Server struct {
	rpcPort        string
	sloggerGrpcUrl string
	grpcServer     *grpc.Server
	servers        []io.Closer
}

func NewServer(rpcPort string, sloggerGrpcUrl string) *Server {
	server := &Server{
		rpcPort:        rpcPort,
		sloggerGrpcUrl: sloggerGrpcUrl,
	}

	return server
}

func (s *Server) Run() error {
	s.grpcServer = grpc.NewServer()

	err := s.registerPingPongServer(s.grpcServer)

	if err != nil {
		return errs.Err(err)
	}

	err = s.registerTraceCollectorServer(s.grpcServer)

	if err != nil {
		return errs.Err(err)
	}

	err = s.registerTraceTransporterServer(s.grpcServer)

	if err != nil {
		return errs.Err(err)
	}

	err = s.registerGrpcManagerServer(s.grpcServer)

	if err != nil {
		return errs.Err(err)
	}

	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		slog.Error("Error listening: ", err.Error())

		return errs.Err(err)
	}

	slog.Info("Listening on " + s.rpcPort)

	err = s.grpcServer.Serve(listener)

	if err != nil {
		slog.Error("Error serving: ", err.Error())
	}

	return errs.Err(err)
}

func (s *Server) Close() error {
	slog.Warn("Closing grpc server...")

	if s.grpcServer == nil {
		return nil
	}

	for _, server := range s.servers {
		err := server.Close()

		if err != nil {
			return errs.Err(err)
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
	server, err := trace_collector.NewServer(s.sloggerGrpcUrl)

	if err != nil {
		slog.Error("Error creating collector: ", err.Error())

		return errs.Err(err)
	}

	trace_collector_gen.RegisterTraceCollectorServer(grpcServer, server)

	s.servers = append(s.servers, server)

	return nil
}

func (s *Server) registerTraceTransporterServer(grpcServer *grpc.Server) error {
	server := trace_transporter.NewServer()

	trace_transporter_gen.RegisterTraceTransporterServer(grpcServer, server)

	s.servers = append(s.servers, server)

	return nil
}

func (s *Server) registerGrpcManagerServer(grpcServer *grpc.Server) error {
	server := grpc_manager.NewServer()

	grpc_manager_gen.RegisterGrpcManagerServer(grpcServer, server)

	s.servers = append(s.servers, server)

	return nil
}
