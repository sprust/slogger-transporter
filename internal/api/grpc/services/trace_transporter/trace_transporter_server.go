package trace_transporter

import (
	"context"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/app"
	"strconv"
)

type Server struct {
	app *app.App
	gen.UnimplementedTraceTransporterServer
}

func NewServer() *Server {
	return &Server{}
}

func (p *Server) Create(ctx context.Context, in *gen.TraceTransporterCreateRequest) (*gen.TraceTransporterResponse, error) {
	go func() {
		slog.Info("grpc:[TraceTransporter.Create]: " + strconv.Itoa(len(in.Payload)))
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (p *Server) Update(ctx context.Context, in *gen.TraceTransporterUpdateRequest) (*gen.TraceTransporterResponse, error) {
	go func() {
		slog.Info("grpc:[TraceTransporter.Update]: " + strconv.Itoa(len(in.Payload)))
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}
