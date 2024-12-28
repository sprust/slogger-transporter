package ping_pong

import (
	"context"
	gen "slogger-transporter/internal/grpc/gen/services/ping_pong_gen"
)

type PingPongServer struct {
	gen.UnimplementedPingPongServer
}

func (p *PingPongServer) Ping(ctx context.Context, in *gen.Request) (*gen.Response, error) {
	return &gen.Response{Message: in.Message}, nil
}
