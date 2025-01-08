package events

import (
	"context"
	"log/slog"
	"slogger-transporter/internal/services/queue_service/objects"
)

type WorkerMessageHandlingFailed struct {
	queueName string
	workerId  int
	message   *objects.Message
	err       error
}

func NewWorkerMessageHandlingFailed(queueName string, workerId int, message *objects.Message, err error) *WorkerMessageHandlingFailed {
	return &WorkerMessageHandlingFailed{
		queueName: queueName,
		workerId:  workerId,
		message:   message,
		err:       err,
	}
}

func (e *WorkerMessageHandlingFailed) Handle(ctx context.Context) error {
	slog.Error(
		joinResult(
			makeQueueName(e.queueName),
			makeWorkerName(e.workerId),
			"message handling failed",
			makeMessageName(e.message),
			e.err.Error(),
		),
	)

	return nil
}
