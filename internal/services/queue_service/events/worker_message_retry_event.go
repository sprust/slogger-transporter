package events

import (
	"context"
	"log/slog"
	"slogger/internal/services/queue_service/objects"
	"strconv"
)

type WorkerMessageRetryEvent struct {
	queueName string
	workerId  int
	message   *objects.Message
}

func NewWorkerMessageRetryEvent(queueName string, workerId int, message *objects.Message) *WorkerMessageRetryEvent {
	return &WorkerMessageRetryEvent{
		queueName: queueName,
		workerId:  workerId,
		message:   message,
	}
}

func (e *WorkerMessageRetryEvent) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			makeMessageName(e.message),
			"retrying "+strconv.Itoa(e.message.Tries),
		),
	)
	return nil
}
