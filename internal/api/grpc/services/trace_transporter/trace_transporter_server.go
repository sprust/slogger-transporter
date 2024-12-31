package trace_transporter

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/trace_transporter_service"
	"time"
)

const waitingWorkersEndingInSeconds = 7

type Server struct {
	app                *app.App
	transporterService *trace_transporter_service.Service
	gen.UnimplementedTraceTransporterServer
	closing      bool
	workersCount int
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

func (s *Server) Create(
	ctx context.Context,
	in *gen.TraceTransporterCreateRequest,
) (*gen.TraceTransporterResponse, error) {
	return s.handle(ctx, func(ctx context.Context) {
		s.transporterService.Create(ctx, in.Payload)
	})
}

func (s *Server) Update(
	ctx context.Context,
	in *gen.TraceTransporterUpdateRequest,
) (*gen.TraceTransporterResponse, error) {
	if s.closing {
		return &gen.TraceTransporterResponse{Success: false}, nil
	}

	s.workersCount += 1

	go func() {
		defer func() {
			s.workersCount -= 1
		}()

		s.transporterService.Update(s.prepareContext(ctx), in.Payload)
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Close() error {
	s.closing = true

	if s.transporterService == nil {
		return nil
	}

	start := time.Now()

	if s.workersCount > 0 {
		slog.Info("Waiting for workers to finish...")

		for s.workersCount > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				slog.Info("Force closing workers...")

				break
			}
		}
	}

	return s.transporterService.Close()
}

func (s *Server) handle(ctx context.Context, callback func(ctx context.Context)) (*gen.TraceTransporterResponse, error) {
	if s.closing {
		return &gen.TraceTransporterResponse{Success: false}, nil
	}

	s.workersCount += 1

	go func(ctx context.Context) {
		defer func() {
			s.workersCount -= 1
		}()

		callback(s.prepareContext(ctx))
	}(ctx)

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) prepareContext(ctx context.Context) context.Context {
	preparedContext := context.WithoutCancel(ctx)

	md, ok := metadata.FromIncomingContext(preparedContext)

	if ok {
		preparedContext = metadata.NewOutgoingContext(preparedContext, md)
	}

	return preparedContext
}
