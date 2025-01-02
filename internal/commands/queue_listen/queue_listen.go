package queue_listen

import (
	"os"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/queue_listen_service"
	"strconv"
)

type QueueListenCommand struct {
}

func (c *QueueListenCommand) Title() string {
	return "Start jobs listening"
}

func (c *QueueListenCommand) Parameters() string {
	return "{no parameters}"
}

func (c *QueueListenCommand) Handle(app *app.App, arguments []string) error {
	queueWorkersNum, err := strconv.Atoi(os.Getenv("RABBITMQ_QUEUE_WORKERS_NUM"))

	if err != nil {
		return err
	}

	sloggerGrpcUrl := os.Getenv("SLOGGER_SERVER_GRPC_URL")

	service, err := queue_listen_service.NewService(
		app,
		&queue_listen_service.RmqParams{
			RmqUser:         os.Getenv("RABBITMQ_USER"),
			RmqPass:         os.Getenv("RABBITMQ_PASSWORD"),
			RmqHost:         os.Getenv("RABBITMQ_HOST"),
			RmqPort:         os.Getenv("RABBITMQ_PORT"),
			QueueName:       os.Getenv("RABBITMQ_QUEUE_NAME"),
			QueueWorkersNum: queueWorkersNum,
		},
		sloggerGrpcUrl,
	)

	if err != nil {
		return err
	}

	app.AddCloseListener(service)

	return service.Listen()
}
