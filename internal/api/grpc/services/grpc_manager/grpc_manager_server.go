package grpc_manager

import (
	"context"
	"log/slog"
	"runtime"
	gen "slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
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
		return nil, errs.Err(err)
	}

	// this is never run
	return &gen.GrpcManagerStopResponse{Success: true, Message: "Grpc server stopped"}, nil
}

func (s *Server) Stat(ctx context.Context, in *gen.GrpcManagerStatRequest) (*gen.GrpcManagerStatResponse, error) {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	return &gen.GrpcManagerStatResponse{
		NumGoroutine:  uint64(runtime.NumGoroutine()),
		AllocMiB:      float32(mem.Alloc / 1024 / 1024),
		TotalAllocMiB: float32(mem.TotalAlloc / 1024 / 1024),
		SysMiB:        float32(mem.Sys / 1024 / 1024),
		NumGC:         uint64(mem.NumGC),
	}, nil
}

func (s *Server) Close() error {
	slog.Warn("Closing grpc manager server...")

	return nil
}
