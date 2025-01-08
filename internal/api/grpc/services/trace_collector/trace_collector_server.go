package trace_collector

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
	"time"
)

type Server struct {
	app    *app.App
	client gen.TraceCollectorClient
	gen.UnimplementedTraceCollectorServer
}

func NewServer(app *app.App, sloggerGrpcUrl string) (*Server, error) {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	client, err := grpc.NewClient(sloggerGrpcUrl, options)

	if err != nil {
		return nil, errs.Err(err)
	}

	return &Server{app: app, client: gen.NewTraceCollectorClient(client)}, nil
}

func (c *Server) Create(ctx context.Context, in *gen.TraceCreateRequest) (*gen.TraceCollectorResponse, error) {
	go func(ctx context.Context, in *gen.TraceCreateRequest, client gen.TraceCollectorClient) {
		start := time.Now()

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		response, err := client.Create(context.WithoutCancel(ctx), in)

		messagePrefix := "grpc[TraceCollector.Create]: " + time.Since(start).String()

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		slog.Info(messagePrefix + ": " + response.Message)
	}(ctx, in, c.client)

	return &gen.TraceCollectorResponse{StatusCode: 200, Message: "ok"}, nil
}

func (c *Server) Update(ctx context.Context, in *gen.TraceUpdateRequest) (*gen.TraceCollectorResponse, error) {
	go func(ctx context.Context, in *gen.TraceUpdateRequest, client gen.TraceCollectorClient) {
		start := time.Now()

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		response, err := client.Update(context.WithoutCancel(ctx), in)

		messagePrefix := "grpc[TraceCollector.Update]: " + time.Since(start).String()

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		slog.Info(messagePrefix + ": " + response.Message)
	}(ctx, in, c.client)

	return &gen.TraceCollectorResponse{StatusCode: 200, Message: "ok"}, nil
}

func (p *Server) Close() error {
	slog.Warn("Closing trace collector server...")

	return nil
}
