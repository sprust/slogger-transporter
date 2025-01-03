package grpc_manager

import (
	"context"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger-transporter/internal/app"
)

type Server struct {
	app *app.App
	gen.UnimplementedGrpcManagerServer
}

func NewServer(app *app.App) *Server {
	return &Server{app: app}
}

func (s *Server) Stop(ctx context.Context, in *gen.GrpcManagerStopRequest) (*gen.GrpcManagerStopResponse, error) {
	err := s.app.Close()

	if err != nil {
		return nil, err
	}

	// this is never run
	return &gen.GrpcManagerStopResponse{Success: true, Message: "Grpc server stopped"}, nil
}

func (s *Server) Close() error {
	slog.Warn("Closing grpc manager server...")

	return nil
}
