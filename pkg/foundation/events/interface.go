package events

import "context"

type EventInterface interface {
	Handle(ctx context.Context) error
}
