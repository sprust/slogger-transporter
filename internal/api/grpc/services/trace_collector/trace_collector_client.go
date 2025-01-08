package trace_collector

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"slogger-transporter/internal/api/grpc/gen/services/trace_collector_gen"
	"slogger-transporter/internal/services/errs"
)

type Client struct {
	grpcClient trace_collector_gen.TraceCollectorClient
	conn       *grpc.ClientConn
}

func NewClient(sloggerUrl string) (*Client, error) {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(sloggerUrl, options)

	if err != nil {
		return nil, errs.Err(err)
	}

	client := trace_collector_gen.NewTraceCollectorClient(conn)

	return &Client{grpcClient: client, conn: conn}, nil
}

func (c *Client) Get() trace_collector_gen.TraceCollectorClient {
	return c.grpcClient
}

func (c *Client) Close() error {
	return c.conn.Close()
}
