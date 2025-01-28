package queue

import (
	"errors"
	"fmt"
	"slogger/pkg/foundation/errs"
	"slogger/pkg/services/queue/objects"
	"sync"
)

var service *Service
var once sync.Once

type Service struct {
	config    objects.RmqConfig
	queues    map[string]objects.QueueInterface
	listeners []*Listener
}

func InitQueue(config objects.RmqConfig, queues map[string]objects.QueueInterface) *Service {
	once.Do(func() {
		service = &Service{
			config: config,
			queues: queues,
		}
	})

	return service
}

func GetQueueService() *Service {
	if service == nil {
		panic("queue is not initialized")
	}

	return service
}

func (q *Service) Listen() error {
	for _, queue := range q.queues {
		listener, err := NewListener(
			q.config,
			NewPublisher(q.config),
			queue,
		)

		if err != nil {
			return err
		}

		q.listeners = append(q.listeners, listener)
	}

	waitGroup := sync.WaitGroup{}

	for _, listener := range q.listeners {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			err := listener.Listen()

			if err != nil {
				panic(err)
			}
		}()
	}

	waitGroup.Wait()

	return nil
}

func (q *Service) Publish(queueName string, payload []byte) error {
	foundQueue, ok := q.queues[queueName]

	if !ok {
		return errs.Err(
			errors.New(fmt.Sprintf("queue %s not found", queueName)),
		)
	}

	settings, err := foundQueue.GetSettings()

	if err != nil {
		return errs.Err(err)
	}

	return NewPublisher(q.config).Publish(settings, payload)
}

func (q *Service) Close() error {
	for _, listener := range q.listeners {
		err := listener.Close()

		if err != nil {
			return errs.Err(err)
		}
	}

	return nil
}
