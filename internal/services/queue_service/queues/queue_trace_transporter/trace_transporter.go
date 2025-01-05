package queue_trace_transporter

import (
	"encoding/json"
	"errors"
	"os"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/queue_service/connections"
	"slogger-transporter/internal/services/queue_service/objects"
	"slogger-transporter/internal/services/trace_transporter_service"
	"strconv"
	"sync"
)

type QueueTraceTransporter struct {
	app         *app.App
	transporter *trace_transporter_service.Service
	settings    *objects.QueueSettings
	mutex       sync.Mutex
}

func NewQueueTraceTransporter(app *app.App) (*QueueTraceTransporter, error) {
	transporter, err := trace_transporter_service.NewService(app)

	if err != nil {
		return nil, err
	}

	return &QueueTraceTransporter{
		app:         app,
		transporter: transporter,
	}, nil
}

func (q *QueueTraceTransporter) GetSettings() (*objects.QueueSettings, error) {
	if q.settings != nil {
		return q.settings, nil
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.settings != nil {
		return q.settings, nil
	}

	queueWorkersNum, err := strconv.Atoi(os.Getenv("TRACE_TRANSPORTER_QUEUE_WORKERS_NUM"))

	if err != nil {
		return nil, err
	}

	queueName := os.Getenv("TRACE_TRANSPORTER_QUEUE_NAME")

	if queueName == "" {
		return nil, errors.New("TRACE_TRANSPORTER_QUEUE_NAME is not set")
	}

	q.settings = &objects.QueueSettings{
		QueueName:       queueName,
		QueueWorkersNum: queueWorkersNum,
	}

	return q.settings, nil
}

// Publish TODO: maybe validate via unmarshal
func (q *QueueTraceTransporter) Publish(payload []byte) error {
	connection := connections.NewConnection(q.app)

	settings, err := q.GetSettings()

	if err != nil {
		return err
	}

	err = connection.Publish(settings.QueueName, payload)

	_ = connection.Close()

	return err
}

func (q *QueueTraceTransporter) Handle(job *objects.Job) error {
	var payload Payload

	err := json.Unmarshal(job.Payload, &payload)

	if err != nil {
		return err
	}

	var creatingTraces []*trace_transporter_service.CreatingTrace
	var updatingTraces []*trace_transporter_service.UpdatingTrace

	for _, message := range payload.Data {
		if message.Action == "c" { // push
			var trace CreatingTrace

			err = json.Unmarshal([]byte(message.Trace), &trace)

			if err != nil {
				return err
			}

			creatingTraces = append(creatingTraces, &trace_transporter_service.CreatingTrace{
				TraceId:       trace.TraceId,
				ParentTraceId: trace.ParentTraceId,
				Type:          trace.Type,
				Status:        trace.Status,
				Tags:          trace.Tags,
				Data:          trace.Data,
				Duration:      trace.Duration,
				Memory:        trace.Memory,
				Cpu:           trace.Cpu,
				LoggedAt:      trace.LoggedAt,
			})
		} else if message.Action == "u" { // stop
			var trace UpdatingTrace

			err = json.Unmarshal([]byte(message.Trace), &trace)

			if err != nil {
				return err
			}

			updatingTraces = append(updatingTraces, &trace_transporter_service.UpdatingTrace{
				TraceId:   trace.TraceId,
				Status:    trace.Status,
				Profiling: trace.Profiling,
				Tags:      trace.Tags,
				Data:      trace.Data,
				Duration:  trace.Duration,
				Memory:    trace.Memory,
				Cpu:       trace.Cpu,
			})
		}
	}

	if len(creatingTraces) > 0 {
		q.transporter.Create(payload.Token, creatingTraces) // TODO
	}

	if len(updatingTraces) > 0 {
		q.transporter.Update(payload.Token, updatingTraces) // TODO
	}

	return nil
}
