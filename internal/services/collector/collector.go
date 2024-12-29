package collector

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"slogger-transporter/internal/app"
	gen "slogger-transporter/internal/services/collector/grpc/gen/services/trace_collector_gen"
	"time"
)

type Collector struct {
	app    *app.App
	client gen.TraceCollectorClient
	gen.UnimplementedTraceCollectorServer
}

func NewCollectorServer(app *app.App, sloggerGrpcUrl string) (*Collector, error) {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	client, err := grpc.NewClient(sloggerGrpcUrl, options)

	if err != nil {
		return nil, err
	}

	return &Collector{app: app, client: gen.NewTraceCollectorClient(client)}, nil
}

func (c *Collector) Create(ctx context.Context, in *gen.TraceCreateRequest) (*gen.TraceCollectorResponse, error) {
	go func(ctx context.Context, in *gen.TraceCreateRequest, client gen.TraceCollectorClient) {
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

	return &gen.TraceCollectorResponse{StatusCode: 200, Message: "ok"}, nil
}

func (c *Collector) Update(ctx context.Context, in *gen.TraceUpdateRequest) (*gen.TraceCollectorResponse, error) {
	go func(ctx context.Context, in *gen.TraceUpdateRequest, client gen.TraceCollectorClient) {
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

	return &gen.TraceCollectorResponse{StatusCode: 200, Message: "ok"}, nil
}
