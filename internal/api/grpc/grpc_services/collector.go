package grpc_services

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/app"
	"time"
)

type Collector struct {
	app    *app.App
	client trace_collector_gen.TraceCollectorClient
	trace_collector_gen.UnimplementedTraceCollectorServer
}

func NewCollectorServer(app *app.App, sloggerGrpcUrl string) (*Collector, error) {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	client, err := grpc.NewClient(sloggerGrpcUrl, options)

	if err != nil {
		return nil, err
	}

	return &Collector{app: app, client: trace_collector_gen.NewTraceCollectorClient(client)}, nil
}

func (c *Collector) Create(ctx context.Context, in *trace_collector_gen.TraceCreateRequest) (*trace_collector_gen.TraceCollectorResponse, error) {
	go func(ctx context.Context, in *trace_collector_gen.TraceCreateRequest, client trace_collector_gen.TraceCollectorClient) {
		start := time.Now()

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		response, err := client.Create(context.WithoutCancel(ctx), in)

		messagePrefix := "grpc[Collector.Create]: " + time.Since(start).String()

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		slog.Info(messagePrefix + ": " + response.Message)
	}(ctx, in, c.client)

	return &trace_collector_gen.TraceCollectorResponse{StatusCode: 200, Message: "ok"}, nil
}

func (c *Collector) Update(ctx context.Context, in *trace_collector_gen.TraceUpdateRequest) (*trace_collector_gen.TraceCollectorResponse, error) {
	go func(ctx context.Context, in *trace_collector_gen.TraceUpdateRequest, client trace_collector_gen.TraceCollectorClient) {
		start := time.Now()

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		response, err := client.Update(context.WithoutCancel(ctx), in)

		messagePrefix := "grpc[Collector.Update]: " + time.Since(start).String()

		if err != nil {
			slog.Error(messagePrefix + ": " + err.Error())

			return
		}

		slog.Info(messagePrefix + ": " + response.Message)
	}(ctx, in, c.client)

	return &trace_collector_gen.TraceCollectorResponse{StatusCode: 200, Message: "ok"}, nil
}
