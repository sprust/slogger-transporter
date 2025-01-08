package events

import (
	"slogger-transporter/internal/services/queue_service/objects"
	"strconv"
	"strings"
)

func joinResult(s ...string) string {
	return strings.Join(s, ": ")
}

func makeQueueName(queueName string) string {
	return "q[" + queueName + "]"
}

func makeWorkerName(workerId int) string {
	return "w[" + strconv.Itoa(workerId) + "]"
}

func makeMessageName(message *objects.Message) string {
	return "msg[" + message.Id + "]"
}
