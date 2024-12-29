package ping_pong

import (
	"context"
	"log/slog"
	gen "slogger-transporter/internal/services/collector/grpc/gen/services/ping_pong_gen"
)

type PingPongServer struct {
	gen.UnimplementedPingPongServer
}

func NewPingPongServer() *PingPongServer {
	return &PingPongServer{}
}

func (p *PingPongServer) Ping(ctx context.Context, in *gen.Request) (*gen.Response, error) {
	go slog.Info("Ping: " + in.Message)

	return &gen.Response{Message: in.Message}, nil
}
