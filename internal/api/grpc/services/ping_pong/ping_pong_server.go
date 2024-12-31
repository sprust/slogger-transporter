package ping_pong

import (
	"context"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/ping_pong_gen"
)

type Server struct {
	gen.UnimplementedPingPongServer
}

func NewServer() *Server {
	return &Server{}
}

func (p *Server) Ping(ctx context.Context, in *gen.PingPongPingRequest) (*gen.PingPongPingResponse, error) {
	go slog.Info("Ping: " + in.Message)

	return &gen.PingPongPingResponse{Message: in.Message}, nil
}

func (p *Server) Close() error {
	return nil
}
