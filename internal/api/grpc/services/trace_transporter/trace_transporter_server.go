package trace_transporter

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_transporter_gen"
	"slogger-transporter/internal/api/grpc/services/trace_collector"
	"slogger-transporter/internal/api/grpc/services/trace_transporter/trace_transporter_parsers"
	"slogger-transporter/internal/app"
	"time"
)

type Server struct {
	app          *app.App
	client       *trace_collector.Client
	parserCreate *trace_transporter_parsers.ParserCreate
	parserUpdate *trace_transporter_parsers.ParserUpdate
	gen.UnimplementedTraceTransporterServer
}

func NewServer(app *app.App, sloggerUrl string) (*Server, error) {
	client, err := trace_collector.NewClient(sloggerUrl)

	if err != nil {
		return nil, err
	}

	return &Server{
		app:          app,
		client:       client,
		parserCreate: trace_transporter_parsers.NewParserCreate(),
		parserUpdate: trace_transporter_parsers.NewParserUpdate(),
	}, nil
}

func (s *Server) Create(ctx context.Context, in *gen.TraceTransporterCreateRequest) (*gen.TraceTransporterResponse, error) {
	go func() {
		messagePrefix := "grpc[TraceTransporter.Create]: "

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		start := time.Now()

		request, err := s.parserCreate.Parse(in.Payload)

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		response, err := s.client.Get().Create(context.WithoutCancel(ctx), request)

		messagePrefix = messagePrefix + time.Since(start).String()

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		slog.Info(messagePrefix + ": " + response.Message)
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Update(ctx context.Context, in *gen.TraceTransporterUpdateRequest) (*gen.TraceTransporterResponse, error) {
	go func() {
		messagePrefix := "grpc[TraceTransporter.Update]"

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		start := time.Now()

		request, err := s.parserUpdate.Parse(in.Payload)

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		response, err := s.client.Get().Update(context.WithoutCancel(ctx), request)

		messagePrefix = messagePrefix + ": " + time.Since(start).String()

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		slog.Info(messagePrefix + ": " + response.Message)
	}()

	return &gen.TraceTransporterResponse{Success: true}, nil
}

func (s *Server) Close() error {
	return s.client.Close()
}
