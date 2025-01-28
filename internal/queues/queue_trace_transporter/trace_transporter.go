package queue_trace_transporter

import (
	"encoding/json"
	"errors"
	"slogger/internal/services/trace_transporter_service"
	"slogger/pkg/foundation/errs"
	"slogger/pkg/services/queue/objects"
	"strings"
	"sync"
)

type QueueTraceTransporter struct {
	queueName       string
	queueWorkersNum int
	transporter     *trace_transporter_service.Service
	settings        *objects.QueueSettings
	mutex           sync.Mutex
}

func NewQueueTraceTransporter(queueName string, queueWorkersNum int) (*QueueTraceTransporter, error) {
	if queueName == "" {
		return nil, errors.New("invalid queue name")
	}

	if queueWorkersNum < 1 {
		return nil, errors.New("invalid queue workers num")
	}

	transporter, err := trace_transporter_service.NewService()

	if err != nil {
		return nil, errs.Err(err)
	}

	return &QueueTraceTransporter{
		queueName:       queueName,
		queueWorkersNum: queueWorkersNum,
		transporter:     transporter,
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

	q.settings = &objects.QueueSettings{
		QueueName:       q.queueName,
		QueueWorkersNum: q.queueWorkersNum,
	}

	return q.settings, nil
}

func (q *QueueTraceTransporter) Handle(job *objects.Job) error {
	var payload Payload

	err := json.Unmarshal(job.Payload, &payload)

	if err != nil {
		return errs.Err(err)
	}

	var creatingTraces []*trace_transporter_service.CreatingTrace
	var updatingTraces []*trace_transporter_service.UpdatingTrace

	for _, action := range payload.Actions {
		if action.Type == "cr" { // create
			var trace CreatingTrace

			err = json.Unmarshal([]byte(action.Data), &trace)

			if err != nil {
				return errs.Err(err)
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
		} else if action.Type == "upd" { // update
			var trace UpdatingTrace

			err = json.Unmarshal([]byte(action.Data), &trace)

			if err != nil {
				return errs.Err(err)
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
		} else {
			return errors.New("unknown action type: " + action.Type)
		}
	}

	var errorsTexts []string

	if len(creatingTraces) > 0 {
		err = q.transporter.Create(payload.Token, creatingTraces)

		if err != nil {
			errorsTexts = append(errorsTexts, "creating: "+err.Error())
		}
	}

	if len(updatingTraces) > 0 {
		err = q.transporter.Update(payload.Token, updatingTraces)

		if err != nil {
			errorsTexts = append(errorsTexts, "updating: "+err.Error())
		}
	}

	if len(errorsTexts) > 0 {
		return errs.Err(errors.New(strings.Join(errorsTexts, ", ")))
	}

	return nil
}
