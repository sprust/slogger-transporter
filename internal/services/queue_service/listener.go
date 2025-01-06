package queue_service

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/queue_service/connections"
	"slogger-transporter/internal/services/queue_service/objects"
	"sync"
	"time"
)

const (
	maxTries                      = 120
	waitingWorkersEndingInSeconds = 10
)

type Listener struct {
	app                *app.App
	queue              objects.QueueInterface
	queueSettings      *objects.QueueSettings
	events             *Events
	rmqParams          *config.RmqParams
	connections        map[int]*connections.Connection
	publisher          *Publisher
	closing            bool
	handlingCount      int
	handlingCountMutex sync.Mutex
}

func NewListener(app *app.App, queue objects.QueueInterface) (*Listener, error) {
	settings, err := queue.GetSettings()

	if err != nil {
		return nil, err
	}

	return &Listener{
		app:           app,
		queue:         queue,
		queueSettings: settings,
		rmqParams:     app.GetConfig().GetRmqConfig(),
		events:        NewEvents(app, settings.QueueName),
		connections:   make(map[int]*connections.Connection),
		publisher:     NewPublisher(app),
	}, nil
}

func (l *Listener) Listen() error {
	err := l.declareQueue()

	if err != nil {
		return err
	}

	waitGroup := sync.WaitGroup{}

	for id := 0; id < l.queueSettings.QueueWorkersNum; id++ {
		waitGroup.Add(1)

		go func() {
			l.startWorker(id)

			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	return nil
}

func (l *Listener) Close() error {
	l.events.Closing()

	l.closing = true

	for workerId, connection := range l.connections {
		err := connection.Close()

		if err != nil {
			l.events.CloseConnectionFailed(workerId, err)
		}
	}

	if l.handlingCount > 0 {
		start := time.Now()

		l.events.JobsFinishWaiting(waitingWorkersEndingInSeconds)

		for l.handlingCount > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				l.events.JobsForceClosing()

				break
			}
		}

		l.events.JobsFinished()
	}

	l.events.Closed()

	return nil
}

func (l *Listener) startWorker(workerId int) {
	isReconnect := false

	for {
		if l.closing {
			break
		}

		if isReconnect {
			l.events.WorkerReconnecting(workerId)

			time.Sleep(1 * time.Second)
		} else {
			isReconnect = true
		}

		connection := connections.NewConnection(l.app)

		l.events.WorkerConnected(workerId)

		l.connections[workerId] = connection

		deliveries, err := connection.Consume(l.queueSettings.QueueName)

		if err != nil {
			l.events.WorkerRegisterConsumerFailed(workerId, err)

			_ = connection.Close()

			continue
		}

		for delivery := range deliveries {
			l.handleDelivery(workerId, delivery)
		}

		_ = connection.Close()
	}
}

func (l *Listener) declareQueue() error {
	connection := connections.NewConnection(l.app)

	err := connection.DeclareQueue(l.queueSettings.QueueName)

	_ = connection.Close()

	return err
}

func (l *Listener) handleDelivery(workerId int, delivery amqp.Delivery) {
	l.incHandlingCount()

	l.events.WorkerDeliveryReceived(workerId, len(delivery.Body))

	var message objects.Message

	err := json.Unmarshal(delivery.Body, &message)

	if err != nil {
		l.events.WorkerMessageUnmarshalFailed(workerId, err)
		l.decrHandlingCount()
		return
	}

	err = l.queue.Handle(&objects.Job{WorkerId: workerId, Payload: []byte(message.Payload)})

	if err == nil {
		l.decrHandlingCount()
		return
	}

	l.events.WorkerMessageHandlingFailed(workerId, &message, err)

	message.Tries += 1

	if message.Tries > maxTries {
		l.events.WorkerRetryingMessageMaxTriesReached(workerId, &message)
		l.decrHandlingCount()
		return
	}

	go func() {
		defer l.decrHandlingCount()

		l.events.WorkerMessageRetry(workerId, &message)

		var payload []byte

		payload, err = json.Marshal(message)

		if err != nil {
			l.events.WorkerRetryingMessageUnmarshalFailed(workerId, err)

			return
		}

		time.Sleep(1 * time.Second)

		err = l.publisher.Publish(
			l.queueSettings.QueueName,
			payload,
		)

		if err != nil {
			l.events.WorkerRetryingMessagePublishFailed(workerId, &message, err)
		} else {
			l.events.WorkerRetryingMessagePublished(workerId, &message)
		}
	}()
}

func (l *Listener) incHandlingCount() {
	l.handlingCountMutex.Lock()
	defer l.handlingCountMutex.Unlock()

	l.handlingCount += 1
}

func (l *Listener) decrHandlingCount() {
	l.handlingCountMutex.Lock()
	defer l.handlingCountMutex.Unlock()

	l.handlingCount -= 1
}
