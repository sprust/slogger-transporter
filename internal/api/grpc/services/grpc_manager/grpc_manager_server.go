package grpc_manager

import (
	"context"
	"log/slog"
	"runtime"
	gen "slogger-transporter/internal/api/grpc/gen/services/grpc_manager_gen"
	"syscall"
)

type Server struct {
	gen.UnimplementedGrpcManagerServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Stop(ctx context.Context, in *gen.GrpcManagerStopRequest) (*gen.GrpcManagerStopResponse, error) {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	if err != nil {
		return nil, err
	}

	return &gen.GrpcManagerStopResponse{Success: true, Message: "stop signal received..."}, nil
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
