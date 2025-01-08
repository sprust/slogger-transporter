package ping_pong

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	gen "slogger-transporter/internal/api/grpc/gen/services/ping_pong_gen"
	"slogger-transporter/internal/services/errs"
)

type Client struct {
	grpcClient gen.PingPongClient
	conn       *grpc.ClientConn
}

func NewClient(grpcUrl string) (*Client, error) {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(grpcUrl, options)

	if err != nil {
		return nil, errs.Err(err)
	}

	client := gen.NewPingPongClient(conn)

	return &Client{grpcClient: client, conn: conn}, nil
}

func (c *Client) Get() gen.PingPongClient {
	return c.grpcClient
}

func (c *Client) Close() error {
	slog.Warn("Closing ping pong client...")

	return errs.Err(c.conn.Close())
}
