package queue

import (
	"slogger/pkg/services/queue/objects"
)

type ServiceProvider struct {
	config objects.RmqConfig
	queues map[string]objects.QueueInterface
}

func NewQueueServiceProvider(
	config objects.RmqConfig,
	queues map[string]objects.QueueInterface,
) *ServiceProvider {
	return &ServiceProvider{
		config: config,
		queues: queues,
	}
}

func (s *ServiceProvider) Register() error {
	InitQueue(s.config, s.queues)

	return nil
}
