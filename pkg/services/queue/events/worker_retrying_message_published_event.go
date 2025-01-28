package events

import (
	"context"
	"log/slog"
	"slogger/pkg/services/queue/objects"
)

type WorkerRetryingMessagePublished struct {
	queueName string
	workerId  int
	message   *objects.Message
}

func NewWorkerRetryingMessagePublishedEvent(queueName string, workerId int, message *objects.Message) *WorkerRetryingMessagePublished {
	return &WorkerRetryingMessagePublished{
		queueName: queueName,
		workerId:  workerId,
		message:   message,
	}
}

func (e *WorkerRetryingMessagePublished) Handle(ctx context.Context) error {
	slog.Info(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			makeMessageName(e.message),
			"retry published",
		),
	)
	return nil
}
