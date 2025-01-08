package events

import (
	"context"
	"log/slog"
)

type JobsFinishedEvent struct {
	queueName string
}

func NewJobsFinishedEvent(queueName string) *JobsFinishedEvent {
	return &JobsFinishedEvent{
		queueName: queueName,
	}
}

func (e *JobsFinishedEvent) Handle(ctx context.Context) error {
	slog.Info(
		joinResult(
			makeQueueName(e.queueName),
			"jobs finished",
		),
	)

	return nil
}
