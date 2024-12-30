package grpc_services

import (
	"context"
	"log/slog"
	"slogger-transporter/internal/api/grpc/gen/services/ping_pong_gen"
)

type PingPongServer struct {
	ping_pong_gen.UnimplementedPingPongServer
}

func NewPingPongServer() *PingPongServer {
	return &PingPongServer{}
}

func (p *PingPongServer) Ping(ctx context.Context, in *ping_pong_gen.Request) (*ping_pong_gen.Response, error) {
	go slog.Info("Ping: " + in.Message)

	return &ping_pong_gen.Response{Message: in.Message}, nil
}
