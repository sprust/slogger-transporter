package connections

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/errs"
	"sync"
)

type Connection struct {
	app        *app.App
	connection *amqp.Connection
	channel    *amqp.Channel
	mutex      sync.Mutex
}

func NewConnection(app *app.App) *Connection {
	return &Connection{
		app: app,
	}
}

func (c *Connection) DeclareQueue(queueName string) error {
	err := c.Init()

	if err != nil {
		return errs.Err(err)
	}

	_, err = c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return errs.Err(err)
}

func (c *Connection) Consume(queueName string) (<-chan amqp.Delivery, error) {
	err := c.Init()

	if err != nil {
		return nil, errs.Err(err)
	}

	deliveries, err := c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, errs.Err(err)
	}

	return deliveries, nil
}

func (c *Connection) Publish(queueName string, payload []byte) error {
	err := c.Init()

	if err != nil {
		return errs.Err(err)
	}

	return c.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

func (c *Connection) Close() error {
	if c.channel != nil {
		_ = c.channel.Close()
	}

	if c.connection != nil {
		_ = c.connection.Close()
	}

	return nil
}

func (c *Connection) Init() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connection != nil && c.channel != nil {
		return nil
	}

	conn, ch, err := c.connect()

	if err != nil {
		return errs.Err(err)
	}

	c.connection = conn
	c.channel = ch

	return nil
}

func (c *Connection) connect() (*amqp.Connection, *amqp.Channel, error) {
	rmqParams := config.GetConfig().GetRmqConfig()

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rmqParams.RmqUser,
		rmqParams.RmqPass,
		rmqParams.RmqHost,
		rmqParams.RmqPort,
	)

	connection, err := amqp.Dial(url)

	if err != nil {
		return nil, nil, errs.Err(err)
	}

	channel, err := connection.Channel()

	if err != nil {
		_ = connection.Close()

		return nil, nil, errs.Err(err)
	}

	return connection, channel, nil
}
