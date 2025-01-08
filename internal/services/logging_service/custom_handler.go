package logging_service

import (
	"context"
	"log/slog"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/services/errs"
	"slogger-transporter/internal/services/logging_service/handlers"
)

type CustomHandler struct {
	levelPolicy    *LevelPolicy
	consoleHandler slog.Handler
	fileHandler    slog.Handler
}

func NewCustomHandler(app *app.App, levelPolicy *LevelPolicy) (*CustomHandler, error) {
	fileHandler, err := handlers.NewFileHandler()

	if err != nil {
		return nil, err
	}

	app.AddLastCloseListener(fileHandler)

	handler := &CustomHandler{
		levelPolicy:    levelPolicy,
		fileHandler:    fileHandler,
		consoleHandler: handlers.NewConsoleHandler(),
	}

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
	return h.levelPolicy.Allowed(level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}
