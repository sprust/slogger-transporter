package trace_transporter_parsers

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
)

type ParserUpdate struct {
}

type UpdatingTrace struct {
	TraceId   string    `json:"traceId"`
	Status    string    `json:"status"`
	Profiling *string   `json:"profiling"`
	Tags      *[]string `json:"tags"`
	Data      *string   `json:"data"`
	Duration  *float64  `json:"duration"`
	Memory    *float64  `json:"memory"`
	Cpu       *float64  `json:"cpu"`
}

type UpdatingTraces struct {
	Traces []string `json:"traces"`
}

func NewParserUpdate() *ParserUpdate {
	return &ParserUpdate{}
}

func (c *ParserUpdate) Parse(payload string) (*gen.TraceUpdateRequest, error) {
	var traces UpdatingTraces

	err := json.Unmarshal([]byte(payload), &traces)

	if err != nil {
		return nil, err
	}

	var request gen.TraceUpdateRequest

	for _, traceString := range traces.Traces {
		var trace UpdatingTrace

		err = json.Unmarshal([]byte(traceString), &trace)

		if err != nil {
			return nil, err
		}

		var tags *gen.TagsObject

		if trace.Tags != nil {
			tags = &gen.TagsObject{Items: *trace.Tags}
		}

		var data *wrappers.StringValue

		if trace.Data != nil {
			data = &wrappers.StringValue{Value: *trace.Data}
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

		request.Traces = append(request.Traces, &gen.TraceUpdateObject{
			TraceId:   trace.TraceId,
			Status:    trace.Status,
			Profiling: nil, // TODO: parse profiling
			Tags:      tags,
			Data:      data,
			Duration:  duration,
			Memory:    memory,
			Cpu:       cpu,
		})
	}

	return &request, nil
}
