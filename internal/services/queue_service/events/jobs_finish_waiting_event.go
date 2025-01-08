package events

import (
	"context"
	"log/slog"
	"strconv"
)

type JobsFinishWaitingEvent struct {
	queueName string
	timeout   int
}

func NewJobsFinishWaitingEvent(queueName string, timeout int) *JobsFinishWaitingEvent {
	return &JobsFinishWaitingEvent{
		queueName: queueName,
		timeout:   timeout,
	}
}

func (e *JobsFinishWaitingEvent) Handle(ctx context.Context) error {
	slog.Info(
		joinResult(
			makeQueueName(e.queueName),
			"waiting for jobs to finish "+strconv.Itoa(e.timeout)+" seconds...",
		),
	)

	return nil
}
