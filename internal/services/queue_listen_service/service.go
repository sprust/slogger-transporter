package queue_listen_service

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/metadata"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/trace_transporter_service"
	"sync"
	"time"
)

const (
	maxTries                      = 120
	delay                         = 1 * time.Second
	waitingWorkersEndingInSeconds = 10
)

type Service struct {
	app                *app.App
	events             *Events
	rmqParams          *RmqParams
	transporterService *trace_transporter_service.Service
	connections        map[int]*amqp.Connection
	closing            bool
	retryingCount      int
	retryingCountMutex sync.Mutex
}

type RmqParams struct {
	RmqUser         string
	RmqPass         string
	RmqHost         string
	RmqPort         string
	QueueName       string
	QueueWorkersNum int
}

type Message struct {
	Token   string `json:"token"`
	Action  string `json:"action"`
	Payload string `json:"payload"`
	Tries   int    `json:"tries"`
}

func NewService(app *app.App, rmqParams *RmqParams, sloggerUrl string) (*Service, error) {
	transporterService, err := trace_transporter_service.NewService(app, sloggerUrl)

	if err != nil {
		return nil, err
	}

	return &Service{
		app:                app,
		rmqParams:          rmqParams,
		events:             NewEvents(app),
		transporterService: transporterService,
		connections:        make(map[int]*amqp.Connection),
	}, nil
}

func (s *Service) Listen() error {
	err := s.declareQueues()

	if err != nil {
		return err
	}

	waitGroup := sync.WaitGroup{}

	for id := 0; id < s.rmqParams.QueueWorkersNum; id++ {
		waitGroup.Add(1)

		go func() {
			s.startWorker(id)

			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	if s.retryingCount > 0 {
		start := time.Now()

		s.events.JobsFinishWaiting(waitingWorkersEndingInSeconds)

		for s.retryingCount > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				s.events.JobsForceClosing()

				break
			}
		}

		s.events.JobsFinished()
	}

	return nil
}

func (s *Service) Close() error {
	s.events.Closing()

	s.closing = true

	for workerId, connection := range s.connections {
		err := connection.Close()

		if err != nil {
			s.events.CloseConnectionFailed(workerId, err)
		}
	}

	return nil
}

func (s *Service) startWorker(workerId int) {
	isReconnect := false

	for {
		if s.closing {
			break
		}

		if isReconnect {
			s.events.WorkerReconnecting(workerId)

			time.Sleep(1 * time.Second)
		} else {
			isReconnect = true
		}

		connection, channel, err := s.connectRMQ()

		if err != nil {
			s.events.WorkerConnectionFailed(workerId, err)

			continue
		}

		s.events.WorkerConnected(workerId)

		s.connections[workerId] = connection

		messages, err := channel.Consume(
			s.rmqParams.QueueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			s.events.WorkerRegisterConsumerFailed(workerId, err)

			_ = connection.Close()

			continue
		}

		for delivery := range messages {
			s.handleDelivery(workerId, channel, delivery)
		}

		_ = connection.Close()
	}
}

func (s *Service) connectRMQ() (*amqp.Connection, *amqp.Channel, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		s.rmqParams.RmqUser,
		s.rmqParams.RmqPass,
		s.rmqParams.RmqHost,
		s.rmqParams.RmqPort,
	)

	connection, err := amqp.Dial(url)

	if err != nil {
		return nil, nil, err
	}

	channel, err := connection.Channel()

	if err != nil {
		_ = connection.Close()

		return nil, nil, err
	}

	return connection, channel, nil
}

func (s *Service) declareQueues() error {
	_, channel, err := s.connectRMQ()

	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		s.rmqParams.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func (s *Service) handleDelivery(workerId int, channel *amqp.Channel, delivery amqp.Delivery) {
	s.events.WorkerMessageReceived(workerId, len(delivery.Body))

	var message Message

	err := json.Unmarshal(delivery.Body, &message)

	if err != nil {
		s.events.WorkerMessageUnmarshalFailed(workerId, err)

		return
	}

	s.events.WorkerMessageUnmarshal(workerId, &message)

	ctx := context.WithoutCancel(s.app.GetContext())

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + message.Token,
	})

	ctx = metadata.NewOutgoingContext(ctx, md)

	var errResult error

	if message.Action == "push" {
		errResult = s.transporterService.Create(ctx, message.Payload)
	} else if message.Action == "stop" {
		errResult = s.transporterService.Update(ctx, message.Payload)
	} else {
		s.events.WorkerMessageUnknownAction(workerId, &message)

		return
	}

	if errResult == nil {
		return
	}

	s.events.WorkerMessageHandlingFailed(workerId, &message, errResult)

	message.Tries += 1

	if message.Tries > maxTries {
		s.events.WorkerRetryingMessageMaxTriesReached(workerId, &message)

		return
	}

	s.incRetryingCount()

	go func() {
		defer func() {
			s.decrRetryingCount()
		}()

		s.events.WorkerMessageRetry(workerId, &message)

		var payload []byte

		payload, err = json.Marshal(message)

		if err != nil {
			s.events.WorkerRetryingMessageUnmarshalFailed(workerId, err)

			return
		}

		time.Sleep(1 * time.Second)

		err = channel.Publish(
			"",
			s.rmqParams.QueueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        payload,
				Expiration:  fmt.Sprintf("%d", delay),
			},
		)

		if err != nil {
			s.events.WorkerRetryingMessagePublishFailed(workerId, &message, err)

			return
		}

		s.events.WorkerRetryingMessagePublished(workerId, &message)
	}()
}

func (s *Service) incRetryingCount() {
	s.retryingCountMutex.Lock()
	defer s.retryingCountMutex.Unlock()

	s.retryingCount += 1
}

func (s *Service) decrRetryingCount() {
	s.retryingCountMutex.Lock()
	defer s.retryingCountMutex.Unlock()

	s.retryingCount -= 1
}
