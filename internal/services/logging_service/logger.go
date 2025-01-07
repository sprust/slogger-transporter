package logging_service

import (
	"context"
	"log/slog"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
	"slogger-transporter/internal/services/logging_service/handlers"
)

type CustomHandler struct {
	level          *slog.Level
	consoleHandler slog.Handler
	fileHandler    slog.Handler
}

func NewCustomHandler(app *app.App, level *slog.Level) (*CustomHandler, error) {
	handler := &CustomHandler{
		level: level,
	}

	fileHandler, err := handlers.NewFileHandler()

	if err != nil {
		return nil, err
	}

	app.AddLastCloseListener(fileHandler)

	handler.fileHandler = fileHandler
	handler.consoleHandler = handlers.NewConsoleHandler()

	return handler, nil
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if err := h.consoleHandler.Handle(ctx, r); err != nil {
		return errs.Err(err)
	}
	if err := h.fileHandler.Handle(ctx, r); err != nil {
		return errs.Err(err)
	}

	return nil
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.fileHandler.Enabled(ctx, level) || h.consoleHandler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		fileHandler:    h.fileHandler.WithAttrs(attrs),
		consoleHandler: h.consoleHandler.WithAttrs(attrs),
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		fileHandler:    h.fileHandler.WithGroup(name),
		consoleHandler: h.consoleHandler.WithGroup(name),
	}
}

func (h *CustomHandler) Close() error {
	return nil
}

func (h *CustomHandler) initConsoleHandler() {
}

func (h *CustomHandler) makeHandlerOptions() *slog.HandlerOptions {
	if h.level == nil {
		return nil
	}

	return &slog.HandlerOptions{
		Level: *h.level,
	}
}
