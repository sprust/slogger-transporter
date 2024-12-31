package trace_transporter

import (
	"context"
	"google.golang.org/grpc/metadata"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/trace_transporter_service"
)

type Server struct {
	app                *app.App
	transporterService *trace_transporter_service.Service
	gen.UnimplementedTraceTransporterServer
}

func NewServer(app *app.App, sloggerUrl string) (*Server, error) {
	transporterService, err := trace_transporter_service.NewService(app, sloggerUrl)

	if err != nil {
		return nil, err
	}

	return &Server{
		app:                app,
		transporterService: transporterService,
	}, nil
}

func (s *Server) Create(ctx context.Context, in *gen.TraceTransporterCreateRequest) (*gen.TraceTransporterResponse, error) {
	go func() {
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		s.transporterService.Create(context.WithoutCancel(ctx), in.Payload)
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Update(ctx context.Context, in *gen.TraceTransporterUpdateRequest) (*gen.TraceTransporterResponse, error) {
	go func() {
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		s.transporterService.Update(context.WithoutCancel(ctx), in.Payload)
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}
