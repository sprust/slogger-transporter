package config

import (
	"os"
	"slogger-transporter/internal/services/errs"
	"strconv"
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

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
	})

	return instance
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

func (c *Config) GetLogLevels() string {
	return os.Getenv("LOG_LEVELS")
}

func (c *Config) GetRpcPort() string {
	return os.Getenv("RPC_PORT")
}

func (c *Config) GetGrpcPort() string {
	return os.Getenv("GRPC_PORT")
}

func (c *Config) GetTraceTransporterQueueName() string {
	return os.Getenv("TRACE_TRANSPORTER_QUEUE_NAME")
}

func (c *Config) GetTraceTransporterQueueWorkersNum() (int, error) {
	workersNum, err := strconv.Atoi(os.Getenv("TRACE_TRANSPORTER_QUEUE_WORKERS_NUM"))

	if err != nil {
		return 0, errs.Err(err)
	}

	return workersNum, nil
}

func (c *Config) GetSloggerGrpcUrl() string {
	return os.Getenv("SLOGGER_SERVER_GRPC_URL")
}
