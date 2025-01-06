package queue_service

import (
	"errors"
	"fmt"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/queue_service/objects"
	"slogger-transporter/internal/services/queue_service/queues/queue_trace_transporter"
)

type Factory struct {
	items map[string]objects.QueueInterface
}

func NewFactory(app *app.App) (*Factory, error) {
	transporter, err := createTransporter(app)

	if err != nil {
		return nil, err
	}

	return &Factory{
		items: map[string]objects.QueueInterface{
			app.GetConfig().GetTraceTransporterQueueName(): transporter,
		},
	}, err
}

func (f *Factory) GetQueue(name string) (objects.QueueInterface, error) {
	if queue, ok := f.items[name]; ok {
		return queue, nil
	}

	return nil, errors.New(fmt.Sprintf("queue %s not found", name))
}

func createTransporter(app *app.App) (*queue_trace_transporter.QueueTraceTransporter, error) {
	queueWorkersNum, err := app.GetConfig().GetTraceTransporterQueueWorkersNum()

	if err != nil {
		return nil, err
	}

	transporter, err := queue_trace_transporter.NewQueueTraceTransporter(
		app,
		app.GetConfig().GetTraceTransporterQueueName(),
		queueWorkersNum,
	)

	if err != nil {
		return nil, err
	}

	return transporter, nil
}
