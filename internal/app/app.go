package app

import (
	"context"
	"io"
	"log/slog"
	"slogger-transporter/internal/config"
)

type App struct {
	ctx            context.Context
	config         *config.Config
	closeListeners []io.Closer
}

func NewApp(ctx context.Context) App {
	app := App{
		ctx:    ctx,
		config: &config.Config{},
	}

	return app
}

func (a *App) GetContext() context.Context {
	return a.ctx
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) Close() error {
	slog.Info("Closing app...")

	for _, listener := range a.closeListeners {
		err := listener.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) AddCloseListener(listener io.Closer) {
	a.closeListeners = append(a.closeListeners, listener)
}
