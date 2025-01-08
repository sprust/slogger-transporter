package events

import (
	"context"
	"log/slog"
	"sync"
)

type Dispatcher struct {
}

var dispatcher *Dispatcher
var once sync.Once

func GetDispatcher() *Dispatcher {
	once.Do(func() {
		dispatcher = &Dispatcher{}
	})

	return dispatcher
}

func (d *Dispatcher) Dispatch(ctx context.Context, event EventInterface) error {
	return event.Handle(ctx)
}

func (d *Dispatcher) DispatchAsync(ctx context.Context, event EventInterface) {
	go func(ctx context.Context, event EventInterface) {
		err := event.Handle(ctx)

		if err != nil {
			slog.Error(err.Error())
		}
	}(ctx, event)
}
