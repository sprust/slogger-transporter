package events

import (
	"context"
	"log/slog"
)

type WorkerRegisterConsumerFailedEvent struct {
	queueName string
	workerId  int
	err       error
}

func NewWorkerRegisterConsumerFailedEvent(queueName string, workerId int, err error) *WorkerRegisterConsumerFailedEvent {
	return &WorkerRegisterConsumerFailedEvent{
		queueName: queueName,
		workerId:  workerId,
		err:       err,
	}
}

func (e *WorkerRegisterConsumerFailedEvent) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"register consumer failed",
			e.err.Error(),
		),
	)

	return nil
}
