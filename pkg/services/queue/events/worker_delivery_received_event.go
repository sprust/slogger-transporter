package events

import (
	"context"
	"log/slog"
	"strconv"
)

type WorkerDeliveryReceivedEvent struct {
	queueName string
	workerId  int
	bodyLen   int
}

func NewWorkerDeliveryReceivedEvent(queueName string, workerId int, bodyLen int) *WorkerDeliveryReceivedEvent {
	return &WorkerDeliveryReceivedEvent{
		queueName: queueName,
		workerId:  workerId,
		bodyLen:   bodyLen,
	}
}

func (e *WorkerDeliveryReceivedEvent) Handle(ctx context.Context) error {
	slog.Info(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"delivery received",
			"len "+strconv.Itoa(e.bodyLen),
		),
	)

	return nil
}
