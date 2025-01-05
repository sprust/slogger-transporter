package config

import (
	"os"
	"sync"
)

type Config struct {
	rmq   *RmqParams
	mutex sync.Mutex
}

type RmqParams struct {
	RmqUser string
	RmqPass string
	RmqHost string
	RmqPort string
}

func (c *Config) GetRmqConfig() *RmqParams {
	if c.rmq != nil {
		return c.rmq
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.rmq != nil {
		return c.rmq
	}

	c.rmq = &RmqParams{
		RmqUser: os.Getenv("RABBITMQ_USER"),
		RmqPass: os.Getenv("RABBITMQ_PASSWORD"),
		RmqHost: os.Getenv("RABBITMQ_HOST"),
		RmqPort: os.Getenv("RABBITMQ_PORT"),
	}

	return c.rmq
}

func (c *Config) GetSloggerGrpcUrl() string {
	return os.Getenv("SLOGGER_SERVER_GRPC_URL")
}
