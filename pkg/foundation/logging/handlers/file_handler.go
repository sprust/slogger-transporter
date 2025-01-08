package handlers

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"slogger-transporter/pkg/foundation/errs"
	"sync"
	"time"
)

const directory = "logs"

type FileHandler struct {
	logFile            *os.File
	initFileMutex      sync.Mutex
	currentLogFileName string
}

func NewFileHandler() (*FileHandler, error) {
	h := &FileHandler{}

	err := h.freshFileHandler()

	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *FileHandler) Handle(ctx context.Context, r slog.Record) error {
	err := h.freshFileHandler()

	if err != nil {
		return err
	}

	msg := makeMessageByRecord(r) + "\n"

	_, err = h.logFile.WriteString(msg)

	return errs.Err(err)
}

func (h *FileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *FileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *FileHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *FileHandler) Close() error {
	slog.Warn("Closing log file")

	h.initFileMutex.Lock()
	defer h.initFileMutex.Unlock()

	if h.logFile != nil {
		return errs.Err(h.logFile.Close())
	}

	return nil
}

func (h *FileHandler) freshFileHandler() error {
	actualLogFileName := h.makeLogFileName()

	if actualLogFileName == h.currentLogFileName {
		return nil
	}

	h.initFileMutex.Lock()
	defer h.initFileMutex.Unlock()

	if actualLogFileName == h.currentLogFileName {
		return nil
	}

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
	h.currentLogFileName = actualLogFileName

	return nil
}

func (h *FileHandler) makeLogFileName() string {
	return time.Now().Format("2006-01-02") + ".log"
}
