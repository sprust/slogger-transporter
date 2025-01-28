package events

import (
	"context"
	"log/slog"
)

type WorkerMessageUnmarshalFailedEvent struct {
	queueName string
	workerId  int
	err       error
}

func NewWorkerMessageUnmarshalFailedEvent(queueName string, workerId int, err error) *WorkerMessageUnmarshalFailedEvent {
	return &WorkerMessageUnmarshalFailedEvent{
		queueName: queueName,
		workerId:  workerId,
		err:       err,
	}
}

func (e *WorkerMessageUnmarshalFailedEvent) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"message unmarshal failed",
			e.err.Error(),
		),
	)

	return nil
}
