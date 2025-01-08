package events

import (
	"context"
	"log/slog"
)

type WorkerConnectedEvent struct {
	queueName string
	workerId  int
}

func NewWorkerConnectedEvent(queueName string, workerId int) *WorkerConnectedEvent {
	return &WorkerConnectedEvent{
		queueName: queueName,
		workerId:  workerId,
	}
}

func (e *WorkerConnectedEvent) Handle(ctx context.Context) error {
	slog.Info(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"connected",
		),
	)

	return nil
}
