package app

import (
	"context"
	"io"
	"log/slog"
)

type App struct {
	ctx            context.Context
	closeListeners []io.Closer
}

func NewApp(ctx context.Context) App {
	app := App{
		ctx: ctx,
	}

	return app
}

func (a *App) GetContext() context.Context {
	return a.ctx
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
