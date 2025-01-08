package events

import (
	"context"
	"log/slog"
)

type WorkerRetryingMessageUnmarshalFailed struct {
	queueName string
	workerId  int
	err       error
}

func NewWorkerRetryingMessageUnmarshalFailedEvent(queueName string, workerId int, err error) *WorkerRetryingMessageUnmarshalFailed {
	return &WorkerRetryingMessageUnmarshalFailed{
		queueName: queueName,
		workerId:  workerId,
		err:       err,
	}
}

func (e *WorkerRetryingMessageUnmarshalFailed) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"unmarshal error at retrying: "+e.err.Error(),
		),
	)
	return nil
}
