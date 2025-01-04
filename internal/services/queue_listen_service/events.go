package queue_listen_service

import (
	"fmt"
	"log/slog"
	"slogger-transporter/internal/app"
	"strconv"
)

type Events struct {
	app *app.App
}

func NewEvents(app *app.App) *Events {
	return &Events{app: app}
}

func (e *Events) JobsFinishWaiting(timeout int) {
	slog.Info("Waiting for jobs to finish " + strconv.Itoa(timeout) + " seconds...")
}

func (e *Events) JobsForceClosing() {
	slog.Info("Force closing jobs...")
}

func (e *Events) JobsFinished() {
	slog.Info("Jobs finished")
}

func (e *Events) Closing() {
	slog.Warn("Closing queue listen service...")
}

func (e *Events) CloseConnectionFailed(workerId int, err error) {
	slog.Error(fmt.Sprintf("Failed to close connection %d: %s", workerId, err))
}

func (e *Events) WorkerReconnecting(workerId int) {
	slog.Error(fmt.Sprintf("Worker %d: reconnect", workerId))
}

func (e *Events) WorkerConnectionFailed(workerId int, err error) {
	slog.Error(fmt.Sprintf("Worker %d: connection error: %s", workerId, err))
}

func (e *Events) WorkerConnected(workerId int) {
	slog.Info(fmt.Sprintf("Worker %d: connected", workerId))
}

func (e *Events) WorkerRegisterConsumerFailed(workerId int, err error) {
	slog.Error(fmt.Sprintf("Worker %d: failed to register a consumer: %s", workerId, err))
}

func (e *Events) WorkerMessageReceived(workerId int, bodyLen int) {
	slog.Info(fmt.Sprintf("Worker %d: received a message: len %d", workerId, bodyLen))
}

func (e *Events) WorkerMessageUnmarshalFailed(workerId int, err error) {
	slog.Error(fmt.Sprintf("Worker %d: unmarshal error: %s", workerId, err.Error()))
}

func (e *Events) WorkerMessageUnmarshal(workerId int, message *Message) {
	slog.Info(fmt.Sprintf("Worker %d: message unmarshal: action[%s] tries[%d]", workerId, message.Action, message.Tries))
}

func (e *Events) WorkerMessageUnknownAction(workerId int, message *Message) {
	slog.Error(fmt.Sprintf("Worker %d: unknown action: %s", workerId, message.Action))
}

func (e *Events) WorkerMessageHandlingFailed(workerId int, message *Message, err error) {
	slog.Error(fmt.Sprintf("Worker %d: message[%s] handling error: %s", workerId, message.Action, err))
}

func (e *Events) WorkerRetryingMessageMaxTriesReached(workerId int, message *Message) {
	slog.Error(fmt.Sprintf("Worker %d: message[%s] retry: max tries reached", workerId, message.Action))
}

func (e *Events) WorkerMessageRetry(workerId int, message *Message) {
	slog.Info(fmt.Sprintf("Worker %d: message[%s] retry: tries %d", workerId, message.Action, message.Tries))
}

func (e *Events) WorkerRetryingMessageUnmarshalFailed(workerId int, err error) {
	slog.Error(fmt.Sprintf("Worker %d: retry: marshal error: %s", workerId, err))
}

func (e *Events) WorkerRetryingMessagePublishFailed(workerId int, message *Message, err error) {
	slog.Error(fmt.Sprintf("Worker %d: message[%s] retry: publish error: %s", workerId, message.Action, err))
}

func (e *Events) WorkerRetryingMessagePublished(workerId int, message *Message) {
	slog.Info(fmt.Sprintf("Worker %d: message[%s] retry: published", workerId, message.Action))
}
