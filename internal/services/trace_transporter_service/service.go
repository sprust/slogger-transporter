package trace_transporter_service

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/api/grpc/services/trace_collector"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
	"strconv"
	"time"
)

type Service struct {
	app *app.App
}

func NewService(app *app.App) (*Service, error) {
	return &Service{
		app: app,
	}, nil
}

func (s *Service) Create(token string, traces []*CreatingTrace) error {
	var request gen.TraceCreateRequest

	for _, trace := range traces {
		loggedAt, parserErr := time.Parse("2006-01-02 15:04:05.000000", trace.LoggedAt)

		if parserErr != nil {
			return parserErr
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

	if len(request.Traces) == 0 {
		return nil
	}

	ctx := s.makeContext(token)

	client, err := s.getClient()

	if err != nil {
		return errs.Err(err)
	}

	response, err := client.Get().Create(ctx, &request)

	_ = client.Close()

	if err != nil {
		return errs.Err(err)
	}

	if response.GetStatusCode() != 200 {
		return errors.New(
			"transporter[create] invalid status code: " + strconv.Itoa(int(response.StatusCode)) + ", message: " + response.Message,
		)
	}

	return nil
}

func (s *Service) Update(token string, traces []*UpdatingTrace) error {
	var request gen.TraceUpdateRequest

	for _, trace := range traces {
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

	client, err := s.getClient()

	if err != nil {
		return errs.Err(err)
	}

	ctx := s.makeContext(token)

	response, err := client.Get().Update(ctx, &request)

	_ = client.Close()

	if err != nil {
		return errs.Err(err)
	}

	if response.GetStatusCode() != 200 {
		return errors.New(
			"transporter[update] invalid status code: " + strconv.Itoa(int(response.StatusCode)) + ", message: " + response.Message,
		)
	}

	return nil
}

func (s *Service) Close() error {
	slog.Warn("Closing trace transporter service...")

	return nil
}

// CreatingTrace
// coroutines leak in google.golang.org/grpc@v1.69.2/internal/grpcsync/callback_serializer.go:88
// at sharing with server of one application
func (s *Service) getClient() (*trace_collector.Client, error) {
	client, err := trace_collector.NewClient(s.app.GetConfig().GetSloggerGrpcUrl())

	if err != nil {
		return nil, errs.Err(err)
	}

	return client, nil
}

func (s *Service) makeContext(token string) context.Context {
	ctx := context.WithoutCancel(s.app.GetContext())

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
