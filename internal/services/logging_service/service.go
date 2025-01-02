package logging_service

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile   *os.File
	mutex     sync.Mutex
	directory = "logs"
)

func Init() error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	if err := rotateLogFile(); err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(1 * time.Hour)

			if err := rotateLogFile(); err != nil {
				slog.Error("Failed to rotate log file: " + err.Error())
			}
		}
	}()

	customHandler := NewCustomHandler(logFile, os.Stdout)

	logger := slog.New(customHandler)

	slog.SetDefault(logger)

	return nil
}

func rotateLogFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	if logFile != nil {
		_ = logFile.Close()
	}

	logFileName := time.Now().Format("2006-01-02") + ".log"

	var err error

	logFile, err = os.OpenFile(filepath.Join(directory, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	return nil
}
