package queue_listen_service

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/trace_transporter_service"
	"sync"
	"time"
)

var maxTries = 120
var delay = 1 * time.Second

type Service struct {
	app                *app.App
	rmqParams          *RmqParams
	transporterService *trace_transporter_service.Service
	closing            bool
	connections        map[int]*amqp.Connection
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

	return nil
}

func (s *Service) Close() error {
	s.closing = true

	for id, connection := range s.connections {
		err := connection.Close()

		if err != nil {
			slog.Error(fmt.Sprintf("Failed to close connection %d: %s", id, err))
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
			slog.Error(fmt.Sprintf("Worker %d: reconnect", workerId))

			time.Sleep(1 * time.Second)
		} else {
			isReconnect = true
		}

		connection, channel, err := s.connectRMQ()

		if err != nil {
			slog.Error(fmt.Sprintf("Worker %d: connection error: %s", workerId, err.Error()))

			continue
		}

		slog.Info(fmt.Sprintf("Worker %d: connected", workerId))

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
			slog.Error(fmt.Sprintf("Worker %d: failed to register a consumer: %s", workerId, err))

			_ = connection.Close()

			continue
		}

		for d := range messages {
			slog.Info(fmt.Sprintf("Worker %d: received a message: len %d", workerId, len(d.Body)))

			var message Message

			err = json.Unmarshal(d.Body, &message)

			if err != nil {
				slog.Error(fmt.Sprintf("Worker %d: unmarshal error: %s", workerId, err.Error()))

				continue
			}

			go func(message *Message) {
				md := metadata.New(map[string]string{
					"authorization": "Bearer " + message.Token,
				})

				ctx := context.WithoutCancel(s.app.GetContext())

				ctx = metadata.NewOutgoingContext(ctx, md)

				var errResult error

				if message.Action == "push" {
					errResult = s.transporterService.Create(ctx, message.Payload)
				} else if message.Action == "stop" {
					errResult = s.transporterService.Update(ctx, message.Payload)
				} else {
					slog.Error(fmt.Sprintf("Worker %d: unknown action: %s", workerId, message.Action))

					return
				}

				if errResult == nil {
					return
				}

				slog.Error(fmt.Sprintf("Worker %d: error: %s", workerId, errResult.Error()))

				message.Tries += 1

				if message.Tries > maxTries {
					slog.Error(fmt.Sprintf("Worker %d: retry: max tries reached", workerId))

					return
				}

				slog.Info(fmt.Sprintf("Worker %d: retry: tries %d", workerId, message.Tries))

				go func(message *Message, channel *amqp.Channel) {
					payload, err := json.Marshal(message)

					if err != nil {
						slog.Error(fmt.Sprintf("Worker %d: retry: marshal error: %s", workerId, err.Error()))

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
						slog.Error(fmt.Sprintf("Retry: publish error: %s", err.Error()))

						return
					}

					slog.Info(fmt.Sprintf("Worker %d: retry: published", workerId))
				}(message, channel)
			}(&message)
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
