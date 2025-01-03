package commands

import "slogger-transporter/internal/app"

type CommandInterface interface {
	Title() string
	Parameters() string
	Handle(app *app.App, arguments []string) error
	Close() error
}
