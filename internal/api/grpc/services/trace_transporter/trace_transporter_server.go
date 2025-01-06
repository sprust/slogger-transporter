package trace_transporter

import (
	"context"
	"errors"
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
	publisher            *queue_service.Publisher
	closing              bool
	requestHandlingCount int
	requestHandlingMutex sync.Mutex
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
	s.incRequestHandlingCount()

	go func() {
		defer s.decrRequestHandlingCount()

		slog.Info("received trace transporter push request: " + strconv.Itoa(len(in.GetPayload())))

		err := s.publisher.Publish(s.app.GetConfig().GetTraceTransporterQueueName(), []byte(in.GetPayload()))

		if err != nil {
			slog.Error("Failed to publish trace transporter payload: " + err.Error())
		}
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Close() error {
	slog.Warn("Closing trace transporter server...")

	s.closing = true

	if s.requestHandlingCount > 0 {
		start := time.Now()

		slog.Info("Waiting for workers to finish " + strconv.Itoa(waitingWorkersEndingInSeconds) + " seconds...")

		for s.requestHandlingCount > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				slog.Info("Force closing workers...")

				break
			}
		}
	}

	if s.requestHandlingCount > 0 {
		return errors.New("workers are still running")
	}

	return nil
}

func (s *Server) incRequestHandlingCount() {
	s.requestHandlingMutex.Lock()
	defer s.requestHandlingMutex.Unlock()

	s.requestHandlingCount += 1
}

func (s *Server) decrRequestHandlingCount() {
	s.requestHandlingMutex.Lock()
	defer s.requestHandlingMutex.Unlock()

	s.requestHandlingCount -= 1
}
