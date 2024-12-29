package app

import (
	"context"
)

type App struct {
	ctx context.Context
}

func NewApp(ctx context.Context) App {
	app := App{
		ctx: ctx,
	}

	return app
}

func (receiver *App) GetContext() context.Context {
	return receiver.ctx
}
