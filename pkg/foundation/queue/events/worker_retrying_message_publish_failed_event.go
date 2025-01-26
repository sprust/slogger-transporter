package events

import (
	"context"
	"log/slog"
	"slogger/pkg/foundation/queue/objects"
)

type WorkerRetryingMessagePublishFailed struct {
	queueName string
	workerId  int
	message   *objects.Message
	err       error
}

func NewWorkerRetryingMessagePublishFailedEvent(queueName string, workerId int, message *objects.Message, err error) *WorkerRetryingMessagePublishFailed {
	return &WorkerRetryingMessagePublishFailed{
		queueName: queueName,
		workerId:  workerId,
		message:   message,
		err:       err,
	}
}

func (e *WorkerRetryingMessagePublishFailed) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			makeMessageName(e.message),
			"publish error at retrying",
			e.err.Error(),
		),
	)
	return nil
}
