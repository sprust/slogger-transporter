package trace_transporter

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/queue_service"
	"strconv"
	"sync"
	"time"
)

const waitingWorkersEndingInSeconds = 10

type Server struct {
	app *app.App
	gen.UnimplementedTraceTransporterServer
	publisher         *queue_service.Publisher
	closing           bool
	workersCount      int
	workersCountMutex sync.Mutex
}

func NewServer(app *app.App) *Server {
	return &Server{
		app:       app,
		publisher: queue_service.NewPublisher(app),
	}
}

func (s *Server) Push(
	ctx context.Context,
	in *gen.TraceTransporterPushRequest,
) (*gen.TraceTransporterResponse, error) {
	err := s.publisher.Publish(queue_service.QueueTraceTransporterName, []byte(in.GetPayload()))

	if err != nil {
		return &gen.TraceTransporterResponse{Success: false}, err
	}

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Close() error {
	slog.Warn("Closing trace transporter server...")

	s.closing = true

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

	if s.workersCount > 0 {
		return errors.New("workers are still running")
	}

	return nil
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
