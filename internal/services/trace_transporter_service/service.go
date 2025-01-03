package trace_transporter_service

import (
	"context"
	"errors"
	"log/slog"
	"slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/api/grpc/services/trace_collector"
	"slogger-transporter/internal/app"
	"strconv"
	"time"
)

type Service struct {
	app          *app.App
	client       *trace_collector.Client
	parserCreate *ParserCreate
	parserUpdate *ParserUpdate
}

func NewService(app *app.App, sloggerUrl string) (*Service, error) {
	client, err := trace_collector.NewClient(sloggerUrl)

	if err != nil {
		return nil, err
	}

	return &Service{
		app:          app,
		client:       client,
		parserCreate: NewParserCreate(),
		parserUpdate: NewParserUpdate(),
	}, nil
}

func (s *Service) Create(ctx context.Context, payload string) error {
	messagePrefix := "grpc[TraceTransporter.Create]: "

	err := s.send(
		ctx,
		messagePrefix,
		payload,
		func(ctx context.Context, payload string) (*trace_collector_gen.TraceCollectorResponse, error) {
			request, err := s.parserCreate.Parse(payload)

			if err != nil {
				return nil, err
			}

			response, err := s.client.Get().Create(ctx, request)

			if err != nil {
				return nil, err
			}

			return response, nil
		},
	)

	return err
}

func (s *Service) Update(ctx context.Context, payload string) error {
	messagePrefix := "grpc[TraceTransporter.Update]: "

	err := s.send(
		ctx,
		messagePrefix,
		payload,
		func(ctx context.Context, payload string) (*trace_collector_gen.TraceCollectorResponse, error) {
			request, err := s.parserUpdate.Parse(payload)

			if err != nil {
				return nil, err
			}

			response, err := s.client.Get().Update(ctx, request)

			if err != nil {
				return nil, err
			}

			return response, nil
		},
	)

	return err
}

func (s *Service) Close() error {
	return s.client.Close()
}

func (s *Service) send(
	ctx context.Context,
	messagePrefix string,
	payload string,
	callback func(ctx context.Context, payload string) (*trace_collector_gen.TraceCollectorResponse, error),
) error {
	start := time.Now()

	response, err := callback(ctx, payload)

	messagePrefix = messagePrefix + time.Since(start).String()

	if err != nil {
		slog.Error(messagePrefix + ": " + err.Error())

		return err
	}

	if response.GetStatusCode() != 200 {
		return errors.New(
			"status code: " + strconv.Itoa(int(response.StatusCode)) + ", message: " + response.Message,
		)
	}

	slog.Info(messagePrefix + ": " + response.Message)

	return nil
}
