package events

import (
	"context"
	"log/slog"
	"slogger-transporter/internal/services/queue_service/objects"
	"strconv"
)

type WorkerRetryingMessageMaxTriesReachedEvent struct {
	queueName string
	workerId  int
	message   *objects.Message
}

func NewWorkerRetryingMessageMaxTriesReachedEvent(queueName string, workerId int, message *objects.Message) *WorkerRetryingMessageMaxTriesReachedEvent {
	return &WorkerRetryingMessageMaxTriesReachedEvent{
		queueName: queueName,
		workerId:  workerId,
		message:   message,
	}
}

func (e *WorkerRetryingMessageMaxTriesReachedEvent) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			makeMessageName(e.message),
			"max tries reached "+strconv.Itoa(e.message.Tries),
		),
	)
	return nil
}
