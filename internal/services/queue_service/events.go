package queue_service

import (
	"fmt"
	"log/slog"
	"slogger-transporter/internal/services/queue_service/objects"
	"strconv"
)

// TODO: refactor

type Events struct {
	prefix string
}

func NewEvents(queueName string) *Events {
	return &Events{
		prefix: fmt.Sprintf("queue[%s]", queueName),
	}
}

func (e *Events) JobsFinishWaiting(timeout int) {
	slog.Info(e.makeView("waiting for jobs to finish " + strconv.Itoa(timeout) + " seconds..."))
}

func (e *Events) JobsForceClosing() {
	slog.Info(e.makeView("force closing jobs..."))
}

func (e *Events) JobsFinished() {
	slog.Info(e.makeView("jobs finished"))
}

func (e *Events) Closing() {
	slog.Warn(e.makeView("closing queue listen service..."))
}

func (e *Events) Closed() {
	slog.Warn(e.makeView("queue listen service closed"))
}

func (e *Events) CloseConnectionFailed(workerId int, err error) {
	slog.Error(e.makeWorkerLogTextError(workerId, "failed to close connection", err))
}

func (e *Events) WorkerReconnecting(workerId int) {
	slog.Error(e.makeLogText(workerId, "reconnecting"))
}

func (e *Events) WorkerConnectionFailed(workerId int, err error) {
	slog.Error(e.makeWorkerLogTextError(workerId, "connection error", err))
}

func (e *Events) WorkerConnected(workerId int) {
	slog.Info(e.makeLogText(workerId, "connected"))
}

func (e *Events) WorkerRegisterConsumerFailed(workerId int, err error) {
	slog.Error(e.makeLogText(workerId, fmt.Sprintf("failed to register a consumer: %s", err)))
}

func (e *Events) WorkerDeliveryReceived(workerId int, bodyLen int) {
	slog.Info(e.makeLogText(workerId, fmt.Sprintf("received a delivery: len %d", bodyLen)))
}

func (e *Events) WorkerMessageUnmarshalFailed(workerId int, err error) {
	slog.Error(e.makeWorkerLogTextError(workerId, "unmarshal error", err))
}

func (e *Events) WorkerMessageHandlingFailed(workerId int, message *objects.Message, err error) {
	slog.Error(e.makeWorkerLogTextError(workerId, e.makeMessageView(message)+": handling error", err))
}

func (e *Events) WorkerRetryingMessageMaxTriesReached(workerId int, message *objects.Message) {
	slog.Error(e.makeLogText(workerId, e.makeMessageView(message)+": max tries reached"))
}

func (e *Events) WorkerMessageRetry(workerId int, message *objects.Message) {
	slog.Info(e.makeLogText(workerId, e.makeMessageView(message)+fmt.Sprintf(": retry %d", message.Tries)))
}

func (e *Events) WorkerRetryingMessageUnmarshalFailed(workerId int, err error) {
	slog.Error(e.makeWorkerLogTextError(workerId, "unmarshal error at retrying", err))
}

func (e *Events) WorkerRetryingMessagePublishFailed(workerId int, message *objects.Message, err error) {
	slog.Error(e.makeWorkerLogTextError(workerId, e.makeMessageView(message)+": publish error at retrying", err))
}

func (e *Events) WorkerRetryingMessagePublished(workerId int, message *objects.Message) {
	slog.Info(e.makeLogText(workerId, e.makeMessageView(message)+fmt.Sprintf(": retry published %d", message.Tries)))
}

func (e *Events) makeView(text string) string {
	return e.prefix + ": " + text
}

func (e *Events) makeWorkerView(workerId int) string {
	return e.prefix + fmt.Sprintf(": worker[%d]", workerId)
}

func (e *Events) makeLogText(workerId int, text string) string {
	return fmt.Sprintf(e.makeWorkerView(workerId) + ": " + text)
}
func (e *Events) makeWorkerLogTextError(workerId int, text string, err error) string {
	return fmt.Sprintf(e.makeWorkerView(workerId) + ": " + text + ": " + err.Error())
}

func (e *Events) makeMessageView(message *objects.Message) string {
	return fmt.Sprintf("message[%s]", message.Id)
}

func (e *Events) makeWorkerMessageTextError(workerId int, message *objects.Message, text string, err error) string {
	return fmt.Sprintf(e.makeWorkerView(workerId) + ": " + text + ": " + err.Error())
}

func (e *Events) makeWorkerMessageText(workerId int, message *objects.Message, text string, a ...any) string {
	return fmt.Sprintf(
		e.makeWorkerView(workerId) + ": " + e.makeMessageView(message) + ": " + fmt.Sprintf(text, a...),
	)
}
