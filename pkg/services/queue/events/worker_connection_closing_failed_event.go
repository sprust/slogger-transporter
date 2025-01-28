package events

import (
	"context"
	"log/slog"
)

type WorkerConnectionClosingFailedEvent struct {
	queueName string
	workerId  int
	err       error
}

func NewWorkerConnectionClosingFailedEvent(queueName string, workerId int, err error) *WorkerConnectionClosingFailedEvent {
	return &WorkerConnectionClosingFailedEvent{
		queueName: queueName,
		workerId:  workerId,
		err:       err,
	}
}

func (e *WorkerConnectionClosingFailedEvent) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"connection closing failed",
			e.err.Error(),
		),
	)

	return nil
}
