package queue_service

import (
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/queue_service/connections"
	"slogger-transporter/internal/services/queue_service/objects"
	atomic2 "slogger-transporter/pkg/foundation/atomic"
	"slogger-transporter/pkg/foundation/errs"
	"sync"
	"time"
)

const (
	maxTries                      = 120
	waitingWorkersEndingInSeconds = 10
)

type Listener struct {
	queue            objects.QueueInterface
	queueSettings    *objects.QueueSettings
	events           *Events
	rmqParams        *config.RmqParams
	connections      map[int]*connections.Connection
	publisher        *Publisher
	connectionsMutex sync.Mutex
	closing          atomic2.Boolean
	closed           atomic2.Boolean
	handlingCount    atomic2.Counter
}

func NewListener(queue objects.QueueInterface) (*Listener, error) {
	settings, err := queue.GetSettings()

	if err != nil {
		return nil, errs.Err(err)
	}

	return &Listener{
		queue:         queue,
		queueSettings: settings,
		rmqParams:     config.GetConfig().GetRmqConfig(),
		events:        NewEvents(settings.QueueName),
		connections:   make(map[int]*connections.Connection),
		publisher:     NewPublisher(),
	}, nil
}

func (l *Listener) Listen() error {
	if l.closing.Get() {
		return errs.Err(errors.New("listener is closing"))
	}

	l.closed.Set(false)

	err := l.declareQueue()

	if err != nil {
		return errs.Err(err)
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

	for !l.closed.Get() {
		// wait for closing
	}

	return nil
}

func (l *Listener) Close() error {
	l.events.Closing()

	l.closing.Set(true)

	for workerId, connection := range l.connections {
		err := connection.Close()

		if err != nil {
			l.events.CloseConnectionFailed(workerId, err)
		}
	}

	if l.handlingCount.Get() > 0 {
		start := time.Now()

		l.events.JobsFinishWaiting(waitingWorkersEndingInSeconds)

		for l.handlingCount.Get() > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				l.events.JobsForceClosing()

				break
			}
		}

		l.events.JobsFinished()
	}

	l.events.Closed()

	l.closing.Set(false)
	l.closed.Set(true)

	return nil
}

func (l *Listener) startWorker(workerId int) {
	isReconnect := false

	for {
		if l.closing.Get() {
			break
		}

		if isReconnect {
			l.events.WorkerReconnecting(workerId)

			time.Sleep(1 * time.Second)
		} else {
			isReconnect = true
		}

		connection := connections.NewConnection()

		l.events.WorkerConnected(workerId)

		l.addConnection(workerId, connection)

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
	connection := connections.NewConnection()

	err := connection.DeclareQueue(l.queueSettings.QueueName)

	_ = connection.Close()

	return err
}

func (l *Listener) addConnection(workerId int, connection *connections.Connection) {
	l.connectionsMutex.Lock()
	defer l.connectionsMutex.Unlock()

	l.connections[workerId] = connection
}

func (l *Listener) handleDelivery(workerId int, delivery amqp.Delivery) {
	l.handlingCount.Increment()

	l.events.WorkerDeliveryReceived(workerId, len(delivery.Body))

	var message objects.Message

	err := json.Unmarshal(delivery.Body, &message)

	if err != nil {
		l.events.WorkerMessageUnmarshalFailed(workerId, err)
		l.handlingCount.Decrement()
		return
	}

	err = l.queue.Handle(&objects.Job{WorkerId: workerId, Payload: []byte(message.Payload)})

	if err == nil {
		l.handlingCount.Decrement()
		return
	}

	l.events.WorkerMessageHandlingFailed(workerId, &message, err)

	message.Tries += 1

	if message.Tries > maxTries {
		l.events.WorkerRetryingMessageMaxTriesReached(workerId, &message)
		l.handlingCount.Decrement()
		return
	}

	go func() {
		defer l.handlingCount.Decrement()

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
