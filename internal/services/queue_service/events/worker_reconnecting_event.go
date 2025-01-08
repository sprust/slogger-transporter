package events

import (
	"context"
	"log/slog"
)

type WorkerReconnectingEvent struct {
	queueName string
	workerId  int
}

func NewWorkerReconnectingEvent(queueName string, workerId int) *WorkerReconnectingEvent {
	return &WorkerReconnectingEvent{
		queueName: queueName,
		workerId:  workerId,
	}
}

func (e *WorkerReconnectingEvent) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"reconnecting...",
		),
	)

	return nil
}
