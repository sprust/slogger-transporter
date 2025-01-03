package trace_transporter

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/trace_transporter_service"
	"strconv"
	"sync"
	"time"
)

const waitingWorkersEndingInSeconds = 10

type Server struct {
	app                *app.App
	transporterService *trace_transporter_service.Service
	gen.UnimplementedTraceTransporterServer
	closing           bool
	workersCount      int
	workersCountMutex sync.Mutex
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
		_ = s.transporterService.Create(ctx, in.Payload)
	})
}

func (s *Server) Update(
	ctx context.Context,
	in *gen.TraceTransporterUpdateRequest,
) (*gen.TraceTransporterResponse, error) {
	return s.handle(ctx, func(ctx context.Context) {
		_ = s.transporterService.Update(ctx, in.Payload)
	})
}

func (s *Server) Close() error {
	s.closing = true

	if s.transporterService == nil {
		return nil
	}

	if s.workersCount > 0 {
		start := time.Now()

		slog.Info("Waiting for workers to finish " + strconv.Itoa(waitingWorkersEndingInSeconds) + " seconds...")

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

	s.incWorkersCount()

	go func(ctx context.Context) {
		defer func() {
			s.decrWorkersCount()
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

func (s *Server) incWorkersCount() {
	s.workersCountMutex.Lock()
	defer s.workersCountMutex.Unlock()

	s.workersCount += 1
}

func (s *Server) decrWorkersCount() {
	s.workersCountMutex.Lock()
	defer s.workersCountMutex.Unlock()

	s.workersCount -= 1
}
