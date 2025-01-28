package events

import (
	"slogger/pkg/services/queue/objects"
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
