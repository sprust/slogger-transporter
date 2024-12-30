package grpc_services

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"slogger-transporter/internal/api/grpc/gen/services/ping_pong_gen"
	"testing"
	"time"
)

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}

	// TODO: validate env variables
}

func TestPingPong_Ping(t *testing.T) {
	grpcPort := os.Getenv("GRPC_PORT")

	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient("localhost:"+grpcPort, options)

	assert.NoError(t, err)

	client := ping_pong_gen.NewPingPongClient(conn)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()

		assert.NoError(t, err)
	}(conn)

	message := time.Now().String()

	response, err := client.Ping(context.TODO(), &ping_pong_gen.Request{Message: message})

	assert.NoError(t, err)

	assert.Equal(t, message, response.Message)
}
