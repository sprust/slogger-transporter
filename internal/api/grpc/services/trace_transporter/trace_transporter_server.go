package trace_transporter

import (
	"context"
	"errors"
	"log/slog"
	gen "slogger/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger/internal/config"
	"slogger/pkg/foundation/atomic"
	"slogger/pkg/foundation/errs"
	"slogger/pkg/services/queue"
	"strconv"
	"time"
)

const waitingWorkersEndingInSeconds = 10

type Server struct {
	gen.UnimplementedTraceTransporterServer
	queueService         *queue.Service
	closing              atomic.Boolean
	requestHandlingCount atomic.Counter
}

func NewServer() *Server {
	return &Server{
		queueService: queue.GetQueueService(),
	}
}

func (s *Server) Push(
	ctx context.Context,
	in *gen.TraceTransporterPushRequest,
) (*gen.TraceTransporterResponse, error) {
	if s.closing.Get() {
		return nil, errs.Err(errors.New("server is closing"))
	}

	s.requestHandlingCount.Increment()

	go func() {
		s.requestHandlingCount.Decrement()

		slog.Info("received trace transporter push request: " + strconv.Itoa(len(in.GetPayload())))

		err := s.queueService.Publish(
			config.GetConfig().GetTraceTransporterQueueName(),
			[]byte(in.GetPayload()),
		)

		if err != nil {
			slog.Error("Failed to publish trace transporter payload: " + err.Error())
		}
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Close() error {
	slog.Warn("Closing trace transporter server...")

	s.closing.Set(true)

	if s.requestHandlingCount.Get() > 0 {
		start := time.Now()

		slog.Info("Waiting for workers to finish " + strconv.Itoa(waitingWorkersEndingInSeconds) + " seconds...")

		for s.requestHandlingCount.Get() > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				slog.Info("Force closing workers...")

				break
			}
		}
	}

	if s.requestHandlingCount.Get() > 0 {
		return errs.Err(errors.New("workers are still running"))
	}

	return nil
}
