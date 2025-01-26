package events

import (
	"context"
	"log/slog"
)

type QueueListeningClosedEvent struct {
	queueName string
}

func NewQueueListeningClosedEvent(queueName string) *QueueListeningClosedEvent {
	return &QueueListeningClosedEvent{
		queueName: queueName,
	}
}

func (e *QueueListeningClosedEvent) Handle(ctx context.Context) error {
	slog.Warn(
		joinResult(
			makeQueueName(e.queueName),
			"queue listen service closed",
		),
	)

	return nil
}
