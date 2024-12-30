package trace_transporter_parsers

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/timestamppb"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"time"
)

type ParserCreate struct {
}

type CreatingTrace struct {
	TraceId       string   `json:"traceId"`
	ParentTraceId *string  `json:"parentTraceId"`
	Type          string   `json:"type"`
	Status        string   `json:"status"`
	Tags          []string `json:"tags"`
	Data          string   `json:"data"`
	Duration      *float64 `json:"duration"`
	Memory        *float64 `json:"memory"`
	Cpu           *float64 `json:"cpu"`
	LoggedAt      string   `json:"loggedAt"`
}

type CreatingTraces struct {
	Traces []string `json:"traces"`
}

func NewParserCreate() *ParserCreate {
	return &ParserCreate{}
}

func (c *ParserCreate) Parse(payload string) (*gen.TraceCreateRequest, error) {
	var traces CreatingTraces

	err := json.Unmarshal([]byte(payload), &traces)

	if err != nil {
		return nil, err
	}

	var request gen.TraceCreateRequest

	for _, traceString := range traces.Traces {
		var trace CreatingTrace

		err = json.Unmarshal([]byte(traceString), &trace)

		if err != nil {
			return nil, err
		}

		loggedAt, parserErr := time.Parse("2006-01-02 15:04:05.000000", trace.LoggedAt)

		if parserErr != nil {
			return nil, parserErr
		}

		var parentTraceId *wrappers.StringValue

		if trace.ParentTraceId != nil {
			parentTraceId = &wrappers.StringValue{Value: *trace.ParentTraceId}
		}

		var duration *wrappers.DoubleValue

		if trace.Duration != nil {
			duration = &wrappers.DoubleValue{Value: *trace.Duration}
		}

		var memory *wrappers.DoubleValue

		if trace.Memory != nil {
			memory = &wrappers.DoubleValue{Value: *trace.Memory}
		}

		var cpu *wrappers.DoubleValue

		if trace.Cpu != nil {
			cpu = &wrappers.DoubleValue{Value: *trace.Cpu}
		}

		request.Traces = append(request.Traces, &gen.TraceCreateObject{
			TraceId:       trace.TraceId,
			ParentTraceId: parentTraceId,
			Type:          trace.Type,
			Status:        trace.Status,
			Tags:          trace.Tags,
			Data:          trace.Data,
			Duration:      duration,
			Memory:        memory,
			Cpu:           cpu,
			LoggedAt:      timestamppb.New(loggedAt),
		})
	}

	return &request, nil
}
