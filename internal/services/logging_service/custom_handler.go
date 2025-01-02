package logging_service

import (
	"context"
	"log/slog"
	"os"
)

type CustomHandler struct {
	fileHandler    slog.Handler
	consoleHandler slog.Handler
}

func NewCustomHandler(file *os.File, console *os.File) *CustomHandler {
	return &CustomHandler{
		fileHandler:    slog.NewJSONHandler(file, nil),
		consoleHandler: slog.NewTextHandler(console, nil),
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if err := h.fileHandler.Handle(ctx, r); err != nil {
		return err
	}
	if err := h.consoleHandler.Handle(ctx, r); err != nil {
		return err
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
