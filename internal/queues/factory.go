package queues

import (
	"errors"
	"fmt"
	"slogger/internal/config"
	"slogger/internal/queues/queue_trace_transporter"
	"slogger/pkg/foundation/errs"
	"slogger/pkg/services/queue/objects"
)

type Factory struct {
	items map[string]objects.QueueInterface
}

func NewFactory() (*Factory, error) {
	transporter, err := createTransporter()

	if err != nil {
		return nil, errs.Err(err)
	}

	return &Factory{
		items: map[string]objects.QueueInterface{
			config.GetConfig().GetTraceTransporterQueueName(): transporter,
		},
	}, nil
}

func (f *Factory) GetQueue(name string) (objects.QueueInterface, error) {
	if queue, ok := f.items[name]; ok {
		return queue, nil
	}

	return nil, errs.Err(errors.New(fmt.Sprintf("queue %s not found", name)))
}

func createTransporter() (*queue_trace_transporter.QueueTraceTransporter, error) {
	queueWorkersNum, err := config.GetConfig().GetTraceTransporterQueueWorkersNum()

	if err != nil {
		return nil, errs.Err(err)
	}

	transporter, err := queue_trace_transporter.NewQueueTraceTransporter(
		config.GetConfig().GetTraceTransporterQueueName(),
		queueWorkersNum,
	)

	if err != nil {
		return nil, errs.Err(err)
	}

	return transporter, nil
}
