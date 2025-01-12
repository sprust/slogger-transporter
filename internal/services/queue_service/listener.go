package queue_service

import (
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"slogger/internal/config"
	"slogger/internal/services/queue_service/connections"
	"slogger/internal/services/queue_service/events"
	"slogger/internal/services/queue_service/objects"
	"slogger/pkg/foundation/atomic"
	"slogger/pkg/foundation/errs"
	foundationEvents "slogger/pkg/foundation/events"
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
	eventsDispatcher *foundationEvents.Dispatcher
	rmqParams        *config.RmqParams
	connections      map[int]*connections.Connection
	publisher        *Publisher
	connectionsMutex sync.Mutex
	closing          atomic.Boolean
	closed           atomic.Boolean
	handlingCount    atomic.Counter
}

func NewListener(queue objects.QueueInterface) (*Listener, error) {
	settings, err := queue.GetSettings()

	if err != nil {
		return nil, errs.Err(err)
	}

	listener := &Listener{
		queue:            queue,
		queueSettings:    settings,
		rmqParams:        config.GetConfig().GetRmqConfig(),
		eventsDispatcher: foundationEvents.GetDispatcher(),
		connections:      make(map[int]*connections.Connection),
		publisher:        NewPublisher(),
	}

	listener.closed.Set(true)
	listener.closing.Set(false)

	return listener, nil
}

func (l *Listener) Listen() error {
	if l.closing.Get() {
		return errs.Err(errors.New("listener is closing"))
	}

	if !l.closed.Get() {
		return errs.Err(errors.New("listener is not closed"))
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
	_ = l.eventsDispatcher.Dispatch(
		context.TODO(),
		events.NewQueueListeningClosingEvent(l.queueSettings.QueueName),
	)

	l.closing.Set(true)

	for workerId, connection := range l.connections {
		err := connection.Close()

		if err != nil {
			_ = l.eventsDispatcher.Dispatch(
				context.TODO(),
				events.NewWorkerConnectionClosingFailedEvent(l.queueSettings.QueueName, workerId, err),
			)
		}
	}

	if l.handlingCount.Get() > 0 {
		start := time.Now()

		_ = l.eventsDispatcher.Dispatch(
			context.TODO(),
			events.NewJobsFinishWaitingEvent(l.queueSettings.QueueName, waitingWorkersEndingInSeconds),
		)

		for l.handlingCount.Get() > 0 {
			time.Sleep(1 * time.Second)

			if time.Now().Sub(start).Seconds() > waitingWorkersEndingInSeconds {
				_ = l.eventsDispatcher.Dispatch(
					context.TODO(),
					events.NewJobsForceClosingEvent(l.queueSettings.QueueName),
				)

				break
			}
		}

		_ = l.eventsDispatcher.Dispatch(
			context.TODO(),
			events.NewJobsFinishedEvent(l.queueSettings.QueueName),
		)
	}

	_ = l.eventsDispatcher.Dispatch(
		context.TODO(),
		events.NewQueueListeningClosedEvent(l.queueSettings.QueueName),
	)

	l.closing.Set(false)
	l.closed.Set(true)

	return nil
}

func (l *Listener) startWorker(workerId int) {
	isReconnect := false

	for {
		if l.closing.Get() || l.closed.Get() {
			break
		}

		if isReconnect {
			_ = l.eventsDispatcher.Dispatch(
				context.TODO(),
				events.NewWorkerReconnectingEvent(l.queueSettings.QueueName, workerId),
			)

			time.Sleep(1 * time.Second)
		} else {
			isReconnect = true
		}

		connection := connections.NewConnection()

		_ = l.eventsDispatcher.Dispatch(
			context.TODO(),
			events.NewWorkerConnectedEvent(l.queueSettings.QueueName, workerId),
		)

		l.addConnection(workerId, connection)

		deliveries, err := connection.Consume(l.queueSettings.QueueName)

		if err != nil {
			_ = l.eventsDispatcher.Dispatch(
				context.TODO(),
				events.NewWorkerRegisterConsumerFailedEvent(l.queueSettings.QueueName, workerId, err),
			)

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

	_ = l.eventsDispatcher.Dispatch(
		context.TODO(),
		events.NewWorkerDeliveryReceivedEvent(l.queueSettings.QueueName, workerId, len(delivery.Body)),
	)

	var message objects.Message

	err := json.Unmarshal(delivery.Body, &message)

	if err != nil {
		_ = l.eventsDispatcher.Dispatch(
			context.TODO(),
			events.NewWorkerMessageUnmarshalFailedEvent(l.queueSettings.QueueName, workerId, err),
		)
		l.handlingCount.Decrement()
		return
	}

	err = l.queue.Handle(&objects.Job{WorkerId: workerId, Payload: []byte(message.Payload)})

	if err == nil {
		l.handlingCount.Decrement()
		return
	}

	_ = l.eventsDispatcher.Dispatch(
		context.TODO(),
		events.NewWorkerMessageHandlingFailed(l.queueSettings.QueueName, workerId, &message, err),
	)

	message.Tries += 1

	if message.Tries > maxTries {
		_ = l.eventsDispatcher.Dispatch(
			context.TODO(),
			events.NewWorkerRetryingMessageMaxTriesReachedEvent(l.queueSettings.QueueName, workerId, &message),
		)
		l.handlingCount.Decrement()
		return
	}

	go func() {
		defer l.handlingCount.Decrement()

		_ = l.eventsDispatcher.Dispatch(
			context.TODO(),
			events.NewWorkerMessageRetryEvent(l.queueSettings.QueueName, workerId, &message),
		)

		var payload []byte

		payload, err = json.Marshal(message)

		if err != nil {
			_ = l.eventsDispatcher.Dispatch(
				context.TODO(),
				events.NewWorkerRetryingMessageUnmarshalFailedEvent(l.queueSettings.QueueName, workerId, err),
			)

			return
		}

		time.Sleep(1 * time.Second)

		err = l.publisher.Publish(
			l.queueSettings.QueueName,
			payload,
		)

		if err == nil {
			_ = l.eventsDispatcher.Dispatch(
				context.TODO(),
				events.NewWorkerRetryingMessagePublishedEvent(l.queueSettings.QueueName, workerId, &message),
			)
		} else {
			_ = l.eventsDispatcher.Dispatch(
				context.TODO(),
				events.NewWorkerRetryingMessagePublishFailedEvent(l.queueSettings.QueueName, workerId, &message, err),
			)
		}
	}()
}
