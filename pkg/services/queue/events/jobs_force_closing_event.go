package events

import (
	"context"
	"log/slog"
)

type JobsForceClosingEvent struct {
	queueName string
}

func NewJobsForceClosingEvent(queueName string) *JobsForceClosingEvent {
	return &JobsForceClosingEvent{
		queueName: queueName,
	}
}

func (e *JobsForceClosingEvent) Handle(ctx context.Context) error {
	slog.Info(
		joinResult(
			makeQueueName(e.queueName),
			"force closing jobs...",
		),
	)

	return nil
}
