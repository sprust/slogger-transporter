package logging_service

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"slogger-transporter/internal/services/errs"
	"sync"
	"time"
)

const directory = "logs"

// TODO: it write new line as \n

type CustomHandler struct {
	level              *slog.Level
	consoleHandler     slog.Handler
	fileHandler        slog.Handler
	logFile            *os.File
	initFileMutex      sync.Mutex
	currentLogFileName string
}

func NewCustomHandler(level *slog.Level) (*CustomHandler, error) {
	handler := &CustomHandler{
		level: level,
	}

	handler.initConsoleHandler()

	err := handler.initFileHandler()

	if err != nil {
		return nil, errs.Err(err)
	}

	return handler, nil
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if err := h.consoleHandler.Handle(ctx, r); err != nil {
		return errs.Err(err)
	}

	if err := h.initFileHandler(); err != nil {
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
	slog.Warn("Closing log file")

	h.initFileMutex.Lock()
	defer h.initFileMutex.Unlock()

	if h.logFile != nil {
		return errs.Err(h.logFile.Close())
	}

	return nil
}

func (h *CustomHandler) initConsoleHandler() {
	h.consoleHandler = slog.NewTextHandler(os.Stdout, h.makeHandlerOptions())
}

func (h *CustomHandler) initFileHandler() error {
	actualLogFileName := h.makeLogFileName()

	if actualLogFileName == h.currentLogFileName {
		return nil
	}

	h.initFileMutex.Lock()
	defer h.initFileMutex.Unlock()

	filePath := filepath.Join(directory, actualLogFileName)

	if h.logFile != nil {
		_ = h.logFile.Close()
	}

	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		slog.Error("Failed to open log file: " + err.Error())

		return errs.Err(err)
	}

	h.logFile = logFile
	h.fileHandler = slog.NewJSONHandler(h.logFile, h.makeHandlerOptions())
	h.currentLogFileName = actualLogFileName

	return nil
}

func (h *CustomHandler) makeHandlerOptions() *slog.HandlerOptions {
	if h.level == nil {
		return nil
	}

	return &slog.HandlerOptions{
		Level: *h.level,
	}
}

func (h *CustomHandler) makeLogFileName() string {
	return time.Now().Format("2006-01-02") + ".log"
}
