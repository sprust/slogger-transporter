package grpc_manager

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	gen "slogger/internal/api/grpc/gen/services/grpc_manager_gen"
	"slogger/pkg/foundation/errs"
)

type Client struct {
	grpcClient gen.GrpcManagerClient
	conn       *grpc.ClientConn
}

func NewClient(grpcUrl string) (*Client, error) {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(grpcUrl, options)

	if err != nil {
		return nil, errs.Err(err)
	}

	client := gen.NewGrpcManagerClient(conn)

	return &Client{grpcClient: client, conn: conn}, nil
}

func (c *Client) Get() gen.GrpcManagerClient {
	return c.grpcClient
}

func (c *Client) Close() error {
	slog.Warn("Closing grpc manager client...")

	return errs.Err(c.conn.Close())
}
