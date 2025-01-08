package events

import (
	"context"
	"log/slog"
)

type QueueListeningClosingEvent struct {
	queueName string
}

func NewQueueListeningClosingEvent(queueName string) *QueueListeningClosingEvent {
	return &QueueListeningClosingEvent{
		queueName: queueName,
	}
}

func (e *QueueListeningClosingEvent) Handle(ctx context.Context) error {
	slog.Warn(
		joinResult(
			makeQueueName(e.queueName),
			"closing queue listen service...",
		),
	)

	return nil
}
